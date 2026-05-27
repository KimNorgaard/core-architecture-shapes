package main

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Shape struct {
	XML   string  `json:"xml"`
	W     float64 `json:"w"`
	H     float64 `json:"h"`
	Title string  `json:"title,omitempty"`
}

func main() {
	inputPattern := flag.String("input", "dist/optimized/*.svg", "Glob pattern")
	outputFile := flag.String("output", "dist/drawio/core-architecture-shapes.xml", "Output file")
	flag.Parse()

	files, _ := filepath.Glob(*inputPattern)
	var library []Shape
	for _, file := range files {
		shape, err := createDrawIOShape(file)
		if err != nil {
			continue
		}
		library = append(library, shape)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	_ = enc.Encode(library)

	output := fmt.Sprintf("<mxlibrary>%s</mxlibrary>", strings.TrimSpace(buf.String()))
	_ = os.MkdirAll(filepath.Dir(*outputFile), 0755)
	_ = os.WriteFile(*outputFile, []byte(output), 0644)
	fmt.Printf("Created library with %d shapes\n", len(library))
}

func createDrawIOShape(filename string) (Shape, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Shape{}, err
	}

	var paths []string
	decoder := xml.NewDecoder(bytes.NewReader(data))
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "path" {
			for _, attr := range se.Attr {
				if attr.Name.Local == "d" {
					paths = append(paths, attr.Value)
				}
			}
		}
	}

	title := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	title = strings.ReplaceAll(title, "&", "&amp;")

	var foreground strings.Builder
	for _, d := range paths {
		stencilPath := convertDToStencil(d)
		if stencilPath != "" {
			foreground.WriteString("<path>")
			foreground.WriteString(stencilPath)
			foreground.WriteString("</path><fillstroke/>")
		}
	}

	connections := `<connections><constraint x="0.5" y="0"/><constraint x="0.5" y="1"/><constraint x="0" y="0.5"/><constraint x="1" y="0.5"/><constraint x="0.25" y="0"/><constraint x="0.75" y="0"/><constraint x="0.25" y="1"/><constraint x="0.75" y="1"/><constraint x="0" y="0.25"/><constraint x="0" y="0.75"/><constraint x="1" y="0.25"/><constraint x="1" y="0.75"/><constraint x="0" y="0"/><constraint x="1" y="0"/><constraint x="0" y="1"/><constraint x="1" y="1"/></connections>`

	stencilXML := fmt.Sprintf(`<shape name="%s" w="48" h="48" strokewidth="inherit" aspect="fixed">%s<background/><foreground>%s</foreground></shape>`,
		title, connections, foreground.String())

	escapedStencil := url.QueryEscape(stencilXML)
	escapedStencil = strings.ReplaceAll(escapedStencil, "+", "%20")

	var sBuf bytes.Buffer
	sw, _ := flate.NewWriter(&sBuf, flate.BestCompression)
	_, _ = sw.Write([]byte(escapedStencil))
	_ = sw.Close()
	encodedStencil := base64.StdEncoding.EncodeToString(sBuf.Bytes())

	style := fmt.Sprintf("shape=stencil(%s);whiteSpace=wrap;html=1;fillColor=#000000;strokeColor=none;verticalLabelPosition=bottom;verticalAlign=top;align=center;", encodedStencil)

	rawXML := fmt.Sprintf(`<mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/><mxCell id="2" parent="1" style="%s" value="" vertex="1"><mxGeometry height="48" width="48" x="0" y="0" as="geometry"/></mxCell></root></mxGraphModel>`,
		style)

	escapedXML := rawXML
	escapedXML = strings.ReplaceAll(escapedXML, "<", "&lt;")
	escapedXML = strings.ReplaceAll(escapedXML, ">", "&gt;")

	return Shape{
		XML:   escapedXML,
		W:     48,
		H:     48,
		Title: title,
	}, nil
}

func convertDToStencil(d string) string {
	re := regexp.MustCompile(`([a-zA-Z])|([-+]?(?:\d*\.\d+|\d+)(?:[eE][-+]?\d+)?)`)
	matches := re.FindAllStringSubmatch(d, -1)
	var out strings.Builder
	var cmd string
	var args []float64
	var curX, curY, startX, startY float64
	var lastControlX, lastControlY float64

	processCmd := func() {
		if cmd == "" {
			return
		}
		c := cmd[0]
		isRel := c >= "a"[0] && c <= "z"[0]
		upperC := strings.ToUpper(cmd)
		switch upperC {
		case "M":
			for i := 0; i+1 < len(args); i += 2 {
				if isRel {
					curX += args[i]
					curY += args[i+1]
				} else {
					curX = args[i]
					curY = args[i+1]
				}
				if i == 0 {
					fmt.Fprintf(&out, `<move x="%.2f" y="%.2f"/>`, curX, curY)
					startX, startY = curX, curY
				} else {
					fmt.Fprintf(&out, `<line x="%.2f" y="%.2f"/>`, curX, curY)
				}
			}
		case "L":
			for i := 0; i+1 < len(args); i += 2 {
				if isRel {
					curX += args[i]
					curY += args[i+1]
				} else {
					curX = args[i]
					curY = args[i+1]
				}
				fmt.Fprintf(&out, `<line x="%.2f" y="%.2f"/>`, curX, curY)
			}
		case "H":
			for i := 0; i < len(args); i++ {
				if isRel {
					curX += args[i]
				} else {
					curX = args[i]
				}
				fmt.Fprintf(&out, `<line x="%.2f" y="%.2f"/>`, curX, curY)
			}
		case "V":
			for i := 0; i < len(args); i++ {
				if isRel {
					curY += args[i]
				} else {
					curY = args[i]
				}
				fmt.Fprintf(&out, `<line x="%.2f" y="%.2f"/>`, curX, curY)
			}
		case "C":
			for i := 0; i+5 < len(args); i += 6 {
				x1, y1 := args[i], args[i+1]
				x2, y2 := args[i+2], args[i+3]
				x, y := args[i+4], args[i+5]
				if isRel {
					x1 += curX
					y1 += curY
					x2 += curX
					y2 += curY
					x += curX
					y += curY
				}
				fmt.Fprintf(&out, `<curve x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" x3="%.2f" y3="%.2f"/>`, x1, y1, x2, y2, x, y)
				curX, curY = x, y
				lastControlX, lastControlY = x2, y2
			}
		case "S":
			for i := 0; i+3 < len(args); i += 4 {
				x2, y2 := args[i], args[i+1]
				x, y := args[i+2], args[i+3]
				if isRel {
					x2 += curX
					y2 += curY
					x += curX
					y += curY
				}
				x1, y1 := 2*curX-lastControlX, 2*curY-lastControlY
				fmt.Fprintf(&out, `<curve x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" x3="%.2f" y3="%.2f"/>`, x1, y1, x2, y2, x, y)
				curX, curY = x, y
				lastControlX, lastControlY = x2, y2
			}
		case "Q":
			for i := 0; i+3 < len(args); i += 4 {
				x1, y1 := args[i], args[i+1]
				x, y := args[i+2], args[i+3]
				if isRel {
					x1 += curX
					y1 += curY
					x += curX
					y += curY
				}
				fmt.Fprintf(&out, `<quad x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f"/>`, x1, y1, x, y)
				curX, curY = x, y
				lastControlX, lastControlY = x1, y1
			}
		case "T":
			for i := 0; i+1 < len(args); i += 2 {
				x, y := args[i], args[i+1]
				if isRel {
					x += curX
					y += curY
				}
				x1, y1 := 2*curX-lastControlX, 2*curY-lastControlY
				fmt.Fprintf(&out, `<quad x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f"/>`, x1, y1, x, y)
				curX, curY = x, y
				lastControlX, lastControlY = x1, y1
			}
		case "A":
			for i := 0; i+6 < len(args); i += 7 {
				rx, ry := args[i], args[i+1]
				xAxisRot := args[i+2]
				largeArcFlag, sweepFlag := args[i+3], args[i+4]
				x, y := args[i+5], args[i+6]
				if isRel {
					x += curX
					y += curY
				}
				fmt.Fprintf(&out, `<arc rx="%.2f" ry="%.2f" x-axis-rotation="%.2f" large-arc-flag="%.0f" sweep-flag="%.0f" x="%.2f" y="%.2f"/>`, rx, ry, xAxisRot, largeArcFlag, sweepFlag, x, y)
				curX, curY = x, y
			}
		case "Z":
			out.WriteString(`<close/>`)
			curX, curY = startX, startY
		}
		if upperC != "C" && upperC != "S" {
			lastControlX, lastControlY = curX, curY
		}
		args = nil
	}

	for _, m := range matches {
		if m[1] != "" {
			processCmd()
			cmd = m[1]
		} else {
			val, _ := strconv.ParseFloat(m[2], 64)
			args = append(args, val)
		}
	}
	processCmd()
	return out.String()
}
