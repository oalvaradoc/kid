package th

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// RegisterDefaultTranslations registers a set of default translations
// for all built in tag's in validator; you may add your own as desired.
func RegisterDefaultTranslations(v *validator.Validate, trans ut.Translator) (err error) {

	translations := []struct {
		tag             string
		translation     string
		override        bool
		customRegisFunc validator.RegisterTranslationsFunc
		customTransFunc validator.TranslationFunc
	}{
		{
			tag:         "required",
			translation: "{0}เป็นฟิลด์บังคับ",
			override:    false,
		},
		{
			tag: "len",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("len-string", "{0}ความยาวต้องเป็น{1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0}อักขระ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("len-number", "{0}ต้องเท่ากับ{1}", false); err != nil {
					return
				}

				if err = ut.Add("len-items", "{0}ต้องเท่ากับ{1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0}สิ่งของ", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("len-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("len-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("len-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("len-items", fe.Field(), c)

				default:
					t, err = ut.T("len-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("คำเตือน: ช่องแปลผิด: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "min",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("min-string", "{0} ต้องมีความยาวอย่างน้อย {1}", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("min-string-character", "{0} อักขระ", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("min-string-character", "{0} อักขระ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("min-number", "{0} ขั้นต่ำสามารถทำได้เพียง {1}", false); err != nil {
					return
				}

				if err = ut.Add("min-items", "{0} ต้องมีอย่างน้อย {1}", false); err != nil {
					return
				}
				//if err = ut.AddCardinal("min-items-item", "{0} รายการ", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("min-items-item", "{0} รายการ", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("min-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("min-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("min-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("min-items", fe.Field(), c)

				default:
					t, err = ut.T("min-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("คำเตือน: ข้อผิดพลาดในช่องการแปล:% s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "max",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("max-string", "{0} ความยาวต้องไม่เกิน {1}", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("max-string-character", "{0}个字符", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("max-string-character", "{0} อักขระ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("max-number", "{0} ต้องน้อยกว่าหรือเท่ากับ {1}", false); err != nil {
					return
				}

				if err = ut.Add("max-items", "{0} สามารถมีได้มากที่สุด {1}", false); err != nil {
					return
				}
				//if err = ut.AddCardinal("max-items-item", "{0} รายการ", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("max-items-item", "{0} รายการ", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("max-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("max-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("max-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("max-items", fe.Field(), c)

				default:
					t, err = ut.T("max-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("คำเตือน: ข้อผิดพลาดในช่องการแปล:% s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eq",
			translation: "{0} ไม่เท่ากับ {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					fmt.Printf("คำเตือน: ข้อผิดพลาดในช่องการแปล:% # v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ne",
			translation: "{0} ต้องไม่เท่ากับ {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					fmt.Printf("คำเตือน: ข้อผิดพลาดในช่องการแปล:% # v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "lt",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("lt-string", "{0} ต้องมีความยาวน้อยกว่า {1}", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("lt-string-character", "{0} อักขระ", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("lt-string-character", "{0} อักขระ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-number", "{0} ต้องน้อยกว่า {1}", false); err != nil {
					return
				}

				if err = ut.Add("lt-items", "{0} ต้องมีน้อยกว่า {1}", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("lt-items-item", "{0} รายการ", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("lt-items-item", "{0} รายการ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-datetime", "{0} ต้องน้อยกว่าวันที่และเวลาปัจจุบัน", false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lt-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lt-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s'ไม่สามารถใช้สำหรับstructชนิดของ", fe.Tag())
					} else {
						t, err = ut.T("lt-datetime", fe.Field())
					}

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("คำเตือน: ข้อผิดพลาดในช่องการแปล:% s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "lte",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("lte-string", "{0} ความยาวต้องไม่เกิน {1}", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("lte-string-character", "{0} character", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("lte-string-character", "{0} อักขระ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-number", "{0} ต้องน้อยกว่าหรือเท่ากับ {1}", false); err != nil {
					return
				}

				if err = ut.Add("lte-items", "{0} สามารถมีได้มากที่สุด {1}", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("lte-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("lte-items-item", "{0} รายการ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-datetime", "{0} ต้องน้อยกว่าหรือเท่ากับวันที่และเวลาปัจจุบัน", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lte-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lte-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s'ไม่สามารถใช้สำหรับstructชนิดของ", fe.Tag())
					} else {
						t, err = ut.T("lte-datetime", fe.Field())
					}

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("คำเตือน: ข้อผิดพลาดในช่องการแปล:% s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gt",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("gt-string", "{0} ต้องมากกว่า {1}", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("gt-string-character", "{0}个字符", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("gt-string-character", "{0} อักขระ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-number", "{0} ต้องมากกว่า {1}", false); err != nil {
					return
				}

				if err = ut.Add("gt-items", "{0} ต้องมากกว่า {1}", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("gt-items-item", "{0} รายการ", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("gt-items-item", "{0} รายการ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-datetime", "{0} ต้องมากกว่าวันที่และเวลาปัจจุบัน", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gt-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gt-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s'ไม่สามารถใช้สำหรับstructชนิดของ", fe.Tag())
					} else {

						t, err = ut.T("gt-datetime", fe.Field())
					}

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("คำเตือน: ช่องการแปลผิด: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gte",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("gte-string", "{0} ต้องมีความยาวอย่างน้อย {1}", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("gte-string-character", "{0} อักขระ", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("gte-string-character", "{0} อักขระ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-number", "{0} ต้องมากกว่าหรือเท่ากับ {1}", false); err != nil {
					return
				}

				if err = ut.Add("gte-items", "{0} ต้องมีอย่างน้อย {1}", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("gte-items-item", "{0}项", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("gte-items-item", "{0} รายการ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-datetime", "{0} ต้องมากกว่าหรือเท่ากับวันที่และเวลาปัจจุบัน", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gte-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gte-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s'ไม่สามารถใช้สำหรับstructชนิดของ", fe.Tag())
					} else {
						t, err = ut.T("gte-datetime", fe.Field())
					}

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("คำเตือน: ช่องแปลผิด: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqfield",
			translation: "{0}ต้องเท่ากับ{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqcsfield",
			translation: "{0}ต้องเท่ากับ{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "necsfield",
			translation: "{0}ไม่สามารถเท่ากับ{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtcsfield",
			translation: "{0}ต้องมากกว่า{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtecsfield",
			translation: "{0}ต้องมากกว่าหรือเท่ากับ{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltcsfield",
			translation: "{0}ต้องน้อยกว่า{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltecsfield",
			translation: "{0}ต้องน้อยกว่าหรือเท่ากับ{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "nefield",
			translation: "{0}ไม่สามารถเท่ากับ{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtfield",
			translation: "{0}ต้องมากกว่า{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtefield",
			translation: "{0}ต้องมากกว่าหรือเท่ากับ{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltfield",
			translation: "{0}ต้องน้อยกว่า{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltefield",
			translation: "{0}ต้องน้อยกว่าหรือเท่ากับ{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "alpha",
			translation: "{0}มีได้เฉพาะตัวอักษร",
			override:    false,
		},
		{
			tag:         "alphanum",
			translation: "{0}มีได้เฉพาะตัวอักษรและตัวเลข",
			override:    false,
		},
		{
			tag:         "numeric",
			translation: "{0}ต้องเป็นค่าที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "number",
			translation: "{0}ต้องเป็นตัวเลขที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "hexadecimal",
			translation: "{0}ต้องเป็นเลขฐานสิบหกที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "hexcolor",
			translation: "{0}ต้องเป็นเลขฐานสิบหกที่ถูกต้องสี",
			override:    false,
		},
		{
			tag:         "rgb",
			translation: "{0}ต้องเป็นสี RGB ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "rgba",
			translation: "{0}ต้องเป็นสี RGBA ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "hsl",
			translation: "{0}ต้องเป็นสี HSL ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "hsla",
			translation: "{0}ต้องเป็นสี HSLA ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "email",
			translation: "{0}ต้องเป็นกล่องจดหมายที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "url",
			translation: "{0}ต้องเป็น URL ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "uri",
			translation: "{0}ต้องเป็น URI ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "base64",
			translation: "{0}ต้องเป็นสตริง Base64 ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "contains",
			translation: "{0}ต้องมีข้อความ'{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "containsany",
			translation: "{0}ต้องมีอักขระต่อไปนี้อย่างน้อยหนึ่งตัว'{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludes",
			translation: "{0}ไม่สามารถมีข้อความได้'{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesall",
			translation: "{0}ต้องไม่มีอักขระใด ๆ ต่อไปนี้'{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesrune",
			translation: "{0}ไม่สามารถมี'{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "isbn",
			translation: "{0}ต้องเป็นหมายเลข ISBN ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "isbn10",
			translation: "{0}ต้องเป็นหมายเลข ISBN-10 ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "isbn13",
			translation: "{0}ต้องเป็นหมายเลข ISBN-13 ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "uuid",
			translation: "{0}ต้องเป็นไฟล์UUID",
			override:    false,
		},
		{
			tag:         "uuid3",
			translation: "{0}ต้องเป็นไฟล์V3 UUID",
			override:    false,
		},
		{
			tag:         "uuid4",
			translation: "{0}ต้องเป็น V4 ที่ถูกต้อง UUID",
			override:    false,
		},
		{
			tag:         "uuid5",
			translation: "{0}ต้องเป็น V5 ที่ถูกต้อง UUID",
			override:    false,
		},
		{
			tag:         "ascii",
			translation: "{0}ต้องมีอักขระ ascii เท่านั้น",
			override:    false,
		},
		{
			tag:         "printascii",
			translation: "{0}ต้องมีอักขระ ascii ที่พิมพ์ได้เท่านั้น",
			override:    false,
		},
		{
			tag:         "multibyte",
			translation: "{0}ต้องมีอักขระหลายไบต์",
			override:    false,
		},
		{
			tag:         "datauri",
			translation: "{0}ต้องมีข้อมูลที่ถูกต้องURI",
			override:    false,
		},
		{
			tag:         "latitude",
			translation: "{0}ต้องมีพิกัดละติจูดที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "longitude",
			translation: "{0}ต้องมีพิกัดลองจิจูดที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "ssn",
			translation: "{0}ต้องเป็นหมายเลขประกันสังคมที่ถูกต้อง(SSN)",
			override:    false,
		},
		{
			tag:         "ipv4",
			translation: "{0}ต้องเป็นที่อยู่ IPv4 ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "ipv6",
			translation: "{0}ต้องเป็นที่อยู่ IPv6 ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "ip",
			translation: "{0}ต้องเป็นที่อยู่ IP ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "cidr",
			translation: "{0} ต้องเป็นเส้นทางระหว่างโดเมนแบบไม่มีคลาส (CIDR) ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "cidrv4",
			translation: "{0} ต้องเป็นเส้นทางระหว่างโดเมน (CIDR) แบบไม่มีคลาสที่ถูกต้องซึ่งมีที่อยู่ IPv4",
			override:    false,
		},
		{
			tag:         "cidrv6",
			translation: "{0} ต้องเป็นเส้นทางระหว่างโดเมน (CIDR) แบบไม่มีคลาสที่ถูกต้องซึ่งมีที่อยู่ IPv6",
			override:    false,
		},
		{
			tag:         "tcp_addr",
			translation: "{0} ต้องเป็นที่อยู่ TCP ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "tcp4_addr",
			translation: "{0} ต้องเป็นที่อยู่ IPv4 TCP ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "tcp6_addr",
			translation: "{0} ต้องเป็นที่อยู่ IPv6 TCP ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "udp_addr",
			translation: "{0} ต้องเป็นที่อยู่ UDP ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "udp4_addr",
			translation: "{0} ต้องเป็นที่อยู่ IPv4 UDP ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "udp6_addr",
			translation: "{0} ต้องเป็นที่อยู่ IPv6 UDP ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "ip_addr",
			translation: "{0} ต้องเป็นที่อยู่ IP ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "ip4_addr",
			translation: "{0} ต้องเป็นที่อยู่ IPv4 ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "ip6_addr",
			translation: "{0} ต้องเป็นที่อยู่ IPv6 ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "unix_addr",
			translation: "{0} ต้องเป็นที่อยู่ UNIX ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "mac",
			translation: "{0} ต้องเป็นที่อยู่ MAC ที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "iscolor",
			translation: "{0} ต้องเป็นสีที่ถูกต้อง",
			override:    false,
		},
		{
			tag:         "oneof",
			translation: "{0} ต้องเป็นหนึ่งใน [{1}]",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				s, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}
				return s
			},
		},
		{
			tag:         "json",
			translation: "{0} ต้องเป็นสตริง JSON",
			override:    false,
		},
		{
			tag:         "lowercase",
			translation: "{0} ต้องเป็นอักษรตัวพิมพ์เล็ก",
			override:    false,
		},
		{
			tag:         "uppercase",
			translation: "{0} ต้องเป็นอักษรตัวพิมพ์ใหญ่",
			override:    false,
		},
		{
			tag:         "datetime",
			translation: "{0} ต้องเป็น {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
	}

	for _, t := range translations {

		if t.customTransFunc != nil && t.customRegisFunc != nil {

			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)

		} else if t.customTransFunc != nil && t.customRegisFunc == nil {

			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)

		} else if t.customTransFunc == nil && t.customRegisFunc != nil {

			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFunc)

		} else {
			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		}

		if err != nil {
			return
		}
	}

	return
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {

	return func(ut ut.Translator) (err error) {

		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return

	}

}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {

	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		log.Printf("คำเตือน: ช่องแปลผิด: %#v", fe)
		return fe.(error).Error()
	}

	return t
}
