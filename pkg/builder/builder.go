package builder

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// XMLData represents the XML structure
type XMLData struct {
	XMLName  xml.Name
	Elements []XMLElement `xml:",any"`
}

// XMLElement represents a single XML element
type XMLElement struct {
	XMLName xml.Name
	Content string `xml:",chardata"`
}

// GenerateStruct generates a Go struct from the XML data
func GenerateStruct(xmlData XMLData) string {
	var sb strings.Builder
	titleCaser := cases.Title(language.English)
	structName := titleCaser.String(xmlData.XMLName.Local)
	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	for _, elem := range xmlData.Elements {
		fieldName := titleCaser.String(elem.XMLName.Local)
		sb.WriteString(fmt.Sprintf("    %s string `xml:\"%s\"`\n", fieldName, elem.XMLName.Local))
	}
	sb.WriteString("}\n")
	return sb.String()
}

func Build_new_struct() {
	// Open XML file
	xmlFile, err := os.Open("input.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return
	}
	defer xmlFile.Close()

	// Read XML file content
	xmlData, err := io.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML file:", err)
		return
	}

	// Unmarshal XML data
	var data XMLData
	err = xml.Unmarshal(xmlData, &data)
	if err != nil {
		fmt.Println("Error unmarshalling XML data:", err)
		return
	}

	// Generate Go struct
	goStruct := GenerateStruct(data)
	fmt.Println(goStruct)
}
