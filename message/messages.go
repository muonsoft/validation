// The source code of the messages is taken from the Symfony Validator component
// See https://github.com/symfony/validator
//
// Copyright (c) 2004-2021 Fabien Potencier
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package message

const (
	NotBlank      = "This value should not be blank."
	Blank         = "This value should be blank."
	NotNil        = "This value should not be nil."
	CountTooFew   = "This collection should contain {{ limit }} element(s) or more."
	CountTooMany  = "This collection should contain {{ limit }} element(s) or less."
	CountExact    = "This collection should contain exactly {{ limit }} element(s)."
	LengthTooFew  = "This value is too short. It should have {{ limit }} character(s) or more."
	LengthTooMany = "This value is too long. It should have {{ limit }} character(s) or less."
	LengthExact   = "This value should have exactly {{ limit }} character(s)."
)
