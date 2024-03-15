package logger

import (
	"testing"
)

func TestMaskString(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          `{"username": "john_doe", "password": "123"} <password> 1234 </password> <credentials> 12345 </credentials> base64data: ZnNkZnNkZnNkZnNkZnNkZnNkZnNkZnNmc2RzZGZzZGZkc2ZzZGZzZGZzZGZmc2QK`,
			expectedOutput: `{"username": "john_doe", "password": "***"} <password> **** </password> <credentials> ***** </credentials> base64data: ****************************************************************`,
		},
		{
			input:          `{"username": "john_doe", "password": "123"} <password> 1234 </password> <credentials> 12345 </credentials> base64data: ZnNkZnNkZnNkZnNkZnNkZnNkZnNkZnNmc2RzZGZzZGZkc2ZzZGZzZGZzZGZmc2QK`,
			expectedOutput: `{"username": "john_doe", "password": "***"} <password> **** </password> <credentials> ***** </credentials> base64data: ****************************************************************`,
		},
		{
			input:          `{"password": "123"} <password>1234</password> <credentials>12345</credentials> base64data: XYZ123==`,
			expectedOutput: `{"password": "***"} <password>****</password> <credentials>*****</credentials> base64data: XYZ123==`,
		},
		{
			input:          `<![CDATA[ <password> 123 </password> ]]>, <![CDATA[ <credentials> user:1234 </credentials> ]]>, base64data: MNO456==`,
			expectedOutput: `<![CDATA[ <password> *** </password> ]]>, <![CDATA[ <credentials> ********* </credentials> ]]>, base64data: MNO456==`,
		},
		{
			input:          `base64data: ABCDEFGH12345==, <password>123</password>, {"password": "1234"}, <credentials>user:1234</credentials>`,
			expectedOutput: `base64data: ABCDEFGH12345==, <password>***</password>, {"password": "****"}, <credentials>*********</credentials>`,
		},
		{
			input:          `nothing`,
			expectedOutput: `nothing`,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := MaskSensitiveData(tc.input)
			if result != tc.expectedOutput {
				t.Errorf("Expected: %s, Got: %s", tc.expectedOutput, result)
			}
		})
	}
}
