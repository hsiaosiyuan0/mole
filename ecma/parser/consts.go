package parser

type ESVersion int

const (
	ES5  ESVersion = 2009        // https://262.ecma-international.org/5.1/
	ES6  ESVersion = iota + 2014 // https://262.ecma-international.org/6.0/
	ES7                          // https://262.ecma-international.org/7.0/
	ES8                          // https://262.ecma-international.org/8.0/
	ES9                          // https://262.ecma-international.org/9.0/
	ES10                         // https://262.ecma-international.org/10.0/
	ES11                         // https://262.ecma-international.org/11.0/
	ES12                         // https://262.ecma-international.org/12.0/
	ES13                         // https://262.ecma-international.org/13.0/
)

// TODO: es version to features group
