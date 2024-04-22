package docx

import (
	"encoding/xml"
	"errors"
)

// ParagraphItem - параграф
type ParagraphItem struct {
	Params ParagraphParams `xml:"pPr"`
	Items  []DocItem
	Attrs  []xml.Attr
}

// ParagraphParams - параметры параграфа
type ParagraphParams struct {
	Style                 *StringValue   `xml:"pStyle,omitempty"`
	Tabs                  *ParagraphTabs `xml:"tabs,omitempty"`
	Spacing               *SpacingValue  `xml:"spacing,omitempty"`
	Jc                    *StringValue   `xml:"jc,omitempty"`
	Bidi                  *IntValue      `xml:"bidi,omitempty"`
	Size                  *StringValue   `xml:"sz"`
	ParagraphRecordParams *RecordParams  `xml:"rPr"`
	Ind                   *ParagraphInd  `xml:"ind"`
}

//                <w:sectPr w:rsidR="00CD6C33" w:rsidSect="004620AD">
//                    <w:headerReference w:type="default" r:id="rId9"/>
//                    <w:footerReference w:type="default" r:id="rId10"/>
//                    <w:type w:val="continuous"/>
//                    <w:pgSz w:w="11910" w:h="16840"/>
//                    <w:pgMar w:top="1418" w:right="539" w:bottom="1701" w:left="522" w:header="425" w:footer="1140"
//                             w:gutter="0"/>
//                    <w:cols w:space="720"/>
//                    <w:docGrid w:linePitch="299"/>
//                </w:sectPr>
//            </w:pPr>

// ParagraphInd = <w:ind w:left="2136" w:right="1209" w:hanging="882" w:firstLine="223"/>
type ParagraphInd struct {
	Left      string `xml:"left,attr"`
	Right     string `xml:"right,attr"`
	Hanging   string `xml:"hanging,attr"`
	FirstLine string `xml:"firstLine,attr"`
}
type ParagraphTab struct {
	Value    string `xml:"val,attr"`
	Position string `xml:"pos,attr"`
}

type ParagraphTabs struct {
	Tab []*ParagraphTab `xml:"tab"`
}

func (item *ParagraphItem) SetAttrs(attrs []xml.Attr) {
	item.Attrs = attrs
}

// Tag - имя тега элемента
func (item *ParagraphItem) Tag() string {
	return "p"
}

// Type - тип элемента
func (item *ParagraphItem) Type() DocItemType {
	return Paragraph
}

// PlainText - текст
func (item *ParagraphItem) PlainText() string {
	var result string
	for _, i := range item.Items {
		tmp := i.PlainText()
		if len(tmp) > 0 {
			result += tmp
		}
	}
	return result
}

// Clone - клонирование
func (item *ParagraphItem) Clone() DocItem {
	result := new(ParagraphItem)
	result.Items = make([]DocItem, 0)
	for _, i := range item.Items {
		if i != nil {
			result.Items = append(result.Items, i.Clone())
		}
	}
	// Клонирование параметров
	if item.Params.Bidi != nil {
		result.Params.Bidi = new(IntValue)
		result.Params.Bidi.Value = item.Params.Bidi.Value
	}
	if item.Params.Tabs != nil {
		result.Params.Tabs = new(ParagraphTabs)
		result.Params.Tabs.Tab = item.Params.Tabs.Tab
	}
	if item.Params.Jc != nil {
		result.Params.Jc = new(StringValue)
		result.Params.Jc.Value = item.Params.Jc.Value
	}
	if item.Params.Spacing != nil {
		result.Params.Spacing = new(SpacingValue)
		result.Params.Spacing.After = item.Params.Spacing.After
		result.Params.Spacing.Before = item.Params.Spacing.Before
		result.Params.Spacing.Line = item.Params.Spacing.Line
		result.Params.Spacing.LineRule = item.Params.Spacing.LineRule
	}
	if item.Params.Style != nil {
		result.Params.Style = new(StringValue)
		result.Params.Style.Value = item.Params.Style.Value
	}
	if item.Params.Size != nil {
		result.Params.Size = new(StringValue)
		result.Params.Size.Value = item.Params.Size.Value
	}
	if item.Params.ParagraphRecordParams != nil {
		result.Params.ParagraphRecordParams = new(RecordParams)
		result.Params.ParagraphRecordParams = item.Params.ParagraphRecordParams
	}
	return result
}

// Декодирование параграфа
func (item *ParagraphItem) decode(decoder *xml.Decoder) error {
	if decoder != nil {
		var end bool
		for !end {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch element := token.(type) {
			case xml.StartElement:
				{
					if element.Name.Local == "pPr" {
						decoder.DecodeElement(&item.Params, &element)
					} else {
						i := decodeItem(&element, decoder)
						if i != nil {
							item.Items = append(item.Items, i)
						}
					}
				}
			case xml.EndElement:
				{
					if element.Name.Local == "p" {
						end = true
					}
				}
			}
		}
		return nil
	}
	return errors.New("Not have decoder")
}

/* КОДИРОВАНИЕ */

// Кодирование параграфа
func (item *ParagraphItem) encode(encoder *xml.Encoder) error {
	if encoder != nil {
		// Начало параграфа
		start := xml.StartElement{Name: xml.Name{Local: item.Tag()}, Attr: item.Attrs}
		if err := encoder.EncodeToken(start); err != nil {
			return err
		}
		// Параметры параграфа
		if err := encoder.EncodeElement(&item.Params, xml.StartElement{Name: xml.Name{Local: "pPr"}}); err != nil {
			return err
		}
		// Кодируем составные элементы
		for _, i := range item.Items {
			if err := i.encode(encoder); err != nil {
				return err
			}
		}
		// Конец параграфа
		if err := encoder.EncodeToken(start.End()); err != nil {
			return err
		}
		return encoder.Flush()
	}
	return errors.New("Not have encoder")
}
