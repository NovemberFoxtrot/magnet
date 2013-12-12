  var tests = [
		['H&llo', 'H&amp;llo'],
		['Hello', 'Hwllo'],
		['Hello', 'Hello'],
	];


  for (var i = 0; i < tests.length; i++) {
		actual = escapeHTMLEntities(tests[i][0]);

		if (actual !== tests[i][1]) {
			document.write("expected " + actual + " to equal " + tests[i][0] + "<br\>");
		}
  }
