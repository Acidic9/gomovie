var END_OF_INPUT = -1;
var arrChrs = new Array("A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","0","1","2","3","4","5","6","7","8","9","+","/");
var reversegetFChars = new Array;
for (var i = 0; i < arrChrs.length; i++) {
	reversegetFChars[arrChrs[i]] = i
}
var getFStr;
var getFCount;
function ntos(e) {
	e = e.toString(16);
	if (e.length == 1)
		e = "0" + e;
	e = "%" + e;
	return unescape(e)
}
function readReversegetF() {
	if (!getFStr)
		return END_OF_INPUT;
	while (true) {
		if (getFCount >= getFStr.length)
			return END_OF_INPUT;
		var e = getFStr.charAt(getFCount);
		getFCount++;
		if (reversegetFChars[e]) {
			return reversegetFChars[e]
		}
		if (e == "A")
			return 0
	}
	return END_OF_INPUT
}

function setgetFStr(e) {
	getFStr = e;
	getFCount = 0
}
function getF(e) {
	setgetFStr(e);
	var t = "";
	var n = new Array(4);
	var r = false;

	while (!r) {
		n[0] = readReversegetF();
		n[1] = readReversegetF();
		if (n[0] == -1 && n[1] == -1) {
			break;
		}
		n[2] = readReversegetF();
		n[3] = readReversegetF();

		t += ntos(n[0] << 2 & 255 | n[1] >> 4);
		if (n[2] != END_OF_INPUT) {
			t += ntos(n[1] << 4 & 255 | n[2] >> 2);
			if (n[3] != END_OF_INPUT) {
				t += ntos(n[2] << 6 & 255 | n[3])
			} else {
				r = true
			}
		} else {
			r = true
		}
	}


	/*while (!r && (n[0] = readReversegetF()) != END_OF_INPUT && (n[1] = readReversegetF()) != END_OF_INPUT) {
		n[2] = readReversegetF();
		n[3] = readReversegetF();
		t += ntos(n[0] << 2 & 255 | n[1] >> 4);
		if (n[2] != END_OF_INPUT) {
			t += ntos(n[1] << 4 & 255 | n[2] >> 2);
			if (n[3] != END_OF_INPUT) {
				t += ntos(n[2] << 6 & 255 | n[3])
			} else {
				r = true
			}
		} else {
			r = true
		}
	}*/
	return t
}
function doit(e) {
	return unescape(getF(getF(e)))
}


console.log(doit('UEdsbWNtRnRaU0J6Y21NOUltaDBkSEE2THk5MGFHVjJhV1JsYjNNdWRIWXZaVzFpWldRdGRqTTJkMjloWkRneGJUUmlMVGN5T0hnME1UQXVhSFJ0YkNJZ2QyVmlhMmwwUVd4c2IzZEdkV3hzVTJOeVpXVnVQU0owY25WbElpQnRiM3BoYkd4dmQyWjFiR3h6WTNKbFpXNDlJblJ5ZFdVaUlHRnNiRzkzWm5Wc2JITmpjbVZsYmowaWRISjFaU0lnWm5KaGJXVmliM0prWlhJOUlqQWlJRzFoY21kcGJuZHBaSFJvUFNJd0lpQnRZWEpuYVc1b1pXbG5hSFE5SWpBaUlITmpjbTlzYkdsdVp6MGlibThpSUhkcFpIUm9QU0kzTWpnaUlHaGxhV2RvZEQwaU5ERXdJajQ4TDJsbWNtRnRaVDQ9'));