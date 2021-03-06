/*
Package validator implements value validations for structs and individual fields based on tags.
It can also handle Cross Field and Cross Struct validation for nested structs and has the ability
to dive into arrays and maps of any type.

Why not a better error message? because this library intends for you to handle your own error messages.

Why should I handle my own errors? Many reasons, for us building an internationalized application
I needed to know the field and what validation failed so that I could provide an error in the users specific language.

	if fieldErr.Field == "Name" {
		switch fieldErr.ErrorTag
		case "required":
			return "Translated string based on field + error"
		default:
		return "Translated string based on field"
	}

Custom Functions

Custom functions can be added

	// Structure
	func customFunc(topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {

		if whatever {
			return false
		}

		return true
	}

	validate.RegisterValidation("custom tag name", customFunc)
	// NOTES: using the same tag name as an existing function
	//        will overwrite the existing one

Cross Field Validation

Cross Field Validation can be implemented, for example Start & End Date range validation

	// NOTE: when calling validate.Struct(val) topStruct will be the top level struct passed
	//       into the function
	//       when calling validate.FieldWithValue(val, field, tag) val will be
	//       whatever you pass, struct, field...
	//       when calling validate.Field(field, tag) val will be nil
	//
	// Because of the specific requirements and field names within each persons project that
	// uses this library it is likely that custom functions will need to be created for your
	// Cross Field Validation needs, however there are some build in Generic Cross Field validations,
	// see Baked In Validators eqfield, nefield, gtfield, gtefield, ltfield, ltefield and Tags below

Multiple Validators

Multiple validators on a field will process in the order defined

	type Test struct {
		Field `validate:"max=10,min=1"`
	}

	// max will be checked then min

Bad Validator definitions are not handled by the library

	type Test struct {
		Field `validate:"min=10,max=0"`
	}

	// this definition of min max will never succeed

Baked In Validators and Tags

NOTE: Baked In Cross field validation only compares fields on the same struct,
if cross field + cross struct validation is needed your own custom validator
should be implemented.

NOTE2: comma is the default separator of validation tags, if you wish to have a comma
included within the parameter i.e. excludesall=, you will need to use the UTF-8 hex
representation 0x2C, which is replaced in the code as a comma, so the above will
become excludesall=0x2C

NOTE3: pipe is the default separator of or validation tags, if you wish to have a pipe
included within the parameter i.e. excludesall=| you will need to use the UTF-8 hex
representation 0x7C, which is replaced in the code as a pipe, so the above will
become excludesall=0x7C

Here is a list of the current built in validators:

	-
		Tells the validation to skip this struct field; this is particularily
		handy in ignoring embedded structs from being validated. (Usage: -)

	|
		This is the 'or' operator allowing multiple validators to be used and
		accepted. (Usage: rbg|rgba) <-- this would allow either rgb or rgba
		colors to be accepted. This can also be combined with 'and' for example
		( Usage: omitempty,rgb|rgba)

	structonly
		When a field that is a nest struct in encountered and contains this flag
		any validation on the nested struct will be run, but none of the nested
		struct fields will be validated. This is usefull if inside of you program
		you know the struct will be valid, but need to verify it has been assigned.
		NOTE: only "required" and "omitempty" can be used on a struct itself.

	omitempty
		Allows conditional validation, for example if a field is not set with
		a value (Determined by the "required" validator) then other validation
		such as min or max won't run, but if a value is set validation will run.
		(Usage: omitempty)

	dive
		This tells the validator to dive into a slice, array or map and validate that
		level of the slice, array or map with the validation tags that follow.
		Multidimensional nesting is also supported, each level you wish to dive will
		require another dive tag. (Usage: dive)
		Example: [][]string with validation tag "gt=0,dive,len=1,dive,required"
		gt=0 will be applied to []
		len=1 will be applied to []string
		required will be applied to string
		Example2: [][]string with validation tag "gt=0,dive,dive,required"
		gt=0 will be applied to []
		[]string will be spared validation
		required will be applied to string

	required
		This validates that the value is not the data types default zero value.
		For numbers ensures value is not zero. For strings ensures value is
		not "". For slices, maps, pointers, interfaces, channels and functions
		ensures the value is not nil.
		(Usage: required)

	len
		For numbers, max will ensure that the value is
		equal to the parameter given. For strings, it checks that
		the string length is exactly that number of characters. For slices,
		arrays, and maps, validates the number of items. (Usage: len=10)

	max
		For numbers, max will ensure that the value is
		less than or equal to the parameter given. For strings, it checks
		that the string length is at most that number of characters. For
		slices, arrays, and maps, validates the number of items. (Usage: max=10)

	min
		For numbers, min will ensure that the value is
		greater or equal to the parameter given. For strings, it checks that
		the string length is at least that number of characters. For slices,
		arrays, and maps, validates the number of items. (Usage: min=10)

	eq
		For strings & numbers, eq will ensure that the value is
		equal to the parameter given. For slices, arrays, and maps,
		validates the number of items. (Usage: eq=10)

	ne
		For strings & numbers, eq will ensure that the value is not
		equal to the parameter given. For slices, arrays, and maps,
		validates the number of items. (Usage: eq=10)

	gt
		For numbers, this will ensure that the value is greater than the
		parameter given. For strings, it checks that the string length
		is greater than that number of characters. For slices, arrays
		and maps it validates the number of items. (Usage: gt=10)
		For time.Time ensures the time value is greater than time.Now.UTC()
		(Usage: gt)

	gte
		Same as 'min' above. Kept both to make terminology with 'len' easier
		(Usage: gte=10)
		For time.Time ensures the time value is greater than or equal to time.Now.UTC()
		(Usage: gte)

	lt
		For numbers, this will ensure that the value is
		less than the parameter given. For strings, it checks
		that the string length is less than that number of characters.
		For slices, arrays, and maps it validates the number of items.
		(Usage: lt=10)
		For time.Time ensures the time value is less than time.Now.UTC()
		(Usage: lt)

	lte
		Same as 'max' above. Kept both to make terminology with 'len' easier
		(Usage: lte=10)
		For time.Time ensures the time value is less than or equal to time.Now.UTC()
		(Usage: lte)

	eqfield
		This will validate the field value against another fields value either within
		a struct or passed in field.
		usage examples are for validation of a password and confirm password:
		Validation on Password field using validate.Struct Usage(eqfield=ConfirmPassword)
		Validating by field validate.FieldWithValue(password, confirmpassword, "eqfield")

	nefield
		This will validate the field value against another fields value either within
		a struct or passed in field.
		usage examples are for ensuring two colors are not the same:
		Validation on Color field using validate.Struct Usage(nefield=Color2)
		Validating by field validate.FieldWithValue(color1, color2, "nefield")

	gtfield
		Only valid for Numbers and time.Time types, this will validate the field value
		against another fields value either within a struct or passed in field.
		usage examples are for validation of a Start and End date:
		Validation on End field using validate.Struct Usage(gtfield=Start)
		Validating by field validate.FieldWithValue(start, end, "gtfield")

	gtefield
		Only valid for Numbers and time.Time types, this will validate the field value
		against another fields value either within a struct or passed in field.
		usage examples are for validation of a Start and End date:
		Validation on End field using validate.Struct Usage(gtefield=Start)
		Validating by field validate.FieldWithValue(start, end, "gtefield")

	ltfield
		Only valid for Numbers and time.Time types, this will validate the field value
		against another fields value either within a struct or passed in field.
		usage examples are for validation of a Start and End date:
		Validation on End field using validate.Struct Usage(ltfield=Start)
		Validating by field validate.FieldWithValue(start, end, "ltfield")

	ltefield
		Only valid for Numbers and time.Time types, this will validate the field value
		against another fields value either within a struct or passed in field.
		usage examples are for validation of a Start and End date:
		Validation on End field using validate.Struct Usage(ltefield=Start)
		Validating by field validate.FieldWithValue(start, end, "ltefield")

	alpha
		This validates that a string value contains alpha characters only
		(Usage: alpha)

	alphanum
		This validates that a string value contains alphanumeric characters only
		(Usage: alphanum)

	numeric
		This validates that a string value contains a basic numeric value.
		basic excludes exponents etc...
		(Usage: numeric)

	hexadecimal
		This validates that a string value contains a valid hexadecimal.
		(Usage: hexadecimal)

	hexcolor
		This validates that a string value contains a valid hex color including
		hashtag (#)
		(Usage: hexcolor)

	rgb
		This validates that a string value contains a valid rgb color
		(Usage: rgb)

	rgba
		This validates that a string value contains a valid rgba color
		(Usage: rgba)

	hsl
		This validates that a string value contains a valid hsl color
		(Usage: hsl)

	hsla
		This validates that a string value contains a valid hsla color
		(Usage: hsla)

	email
		This validates that a string value contains a valid email
		This may not conform to all possibilities of any rfc standard, but neither
		does any email provider accept all posibilities...
		(Usage: email)

	url
		This validates that a string value contains a valid url
		This will accept any url the golang request uri accepts but must contain
		a schema for example http:// or rtmp://
		(Usage: url)

	uri
		This validates that a string value contains a valid uri
		This will accept any uri the golang request uri accepts (Usage: uri)

	base64
		This validates that a string value contains a valid base64 value.
		Although an empty string is valid base64 this will report an empty string
		as an error, if you wish to accept an empty string as valid you can use
		this with the omitempty tag. (Usage: base64)

	contains
		This validates that a string value contains the substring value.
		(Usage: contains=@)

	containsany
		This validates that a string value contains any Unicode code points
		in the substring value. (Usage: containsany=!@#?)

	containsrune
		This validates that a string value contains the supplied rune value.
		(Usage: containsrune=@)

	excludes
		This validates that a string value does not contain the substring value.
		(Usage: excludes=@)

	excludesall
		This validates that a string value does not contain any Unicode code
		points in the substring value. (Usage: excludesall=!@#?)

	excludesrune
		This validates that a string value does not contain the supplied rune value.
		(Usage: excludesrune=@)

	isbn
		This validates that a string value contains a valid isbn10 or isbn13 value.
		(Usage: isbn)

	isbn10
		This validates that a string value contains a valid isbn10 value.
		(Usage: isbn10)

	isbn13
		This validates that a string value contains a valid isbn13 value.
		(Usage: isbn13)

	uuid
		This validates that a string value contains a valid UUID.
		(Usage: uuid)

	uuid3
		This validates that a string value contains a valid version 3 UUID.
		(Usage: uuid3)

	uuid4
		This validates that a string value contains a valid version 4 UUID.
		(Usage: uuid4)

	uuid5
		This validates that a string value contains a valid version 5 UUID.
		(Usage: uuid5)

	ascii
		This validates that a string value contains only ASCII characters.
		NOTE: if the string is blank, this validates as true.
		(Usage: ascii)

	asciiprint
		This validates that a string value contains only printable ASCII characters.
		NOTE: if the string is blank, this validates as true.
		(Usage: asciiprint)

	multibyte
		This validates that a string value contains one or more multibyte characters.
		NOTE: if the string is blank, this validates as true.
		(Usage: multibyte)

	datauri
		This validates that a string value contains a valid DataURI.
		NOTE: this will also validate that the data portion is valid base64
		(Usage: datauri)

	latitude
		This validates that a string value contains a valid latitude.
		(Usage: latitude)

	longitude
		This validates that a string value contains a valid longitude.
		(Usage: longitude)

	ssn
		This validates that a string value contains a valid U.S. Social Security Number.
		(Usage: ssn)

	ip
		This validates that a string value contains a valid IP Adress.
		(Usage: ip)

	ipv4
		This validates that a string value contains a valid v4 IP Adress.
		(Usage: ipv4)

	ipv6
		This validates that a string value contains a valid v6 IP Adress.
		(Usage: ipv6)

	mac
		This validates that a string value contains a valid MAC Adress defined
		by go's ParseMAC accepted formats and types see:
		http://golang.org/src/net/mac.go?s=866:918#L29
		(Usage: mac)

Validator notes:

	regex
		a regex validator won't be added because commas and = signs can be part of
		a regex which conflict with the validation definitions, although workarounds
		can be made, they take away from using pure regex's. Furthermore it's quick
		and dirty but the regex's become harder to maintain and are not reusable, so
		it's as much a programming philosiphy as anything.

		In place of this new validator functions should be created; a regex can be
		used within the validator function and even be precompiled for better efficiency
		within regexes.go.

		And the best reason, you can submit a pull request and we can keep on adding to the
		validation library of this package!

Panics

This package panics when bad input is provided, this is by design, bad code like that should not make it to production.

	type Test struct {
		TestField string `validate:"nonexistantfunction=1"`
	}

	t := &Test{
		TestField: "Test"
	}

	validate.Struct(t) // this will panic
*/
package validator
