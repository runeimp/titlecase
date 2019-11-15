TitleCase
=========

Based on Stuart Colville's `titlecase.py` now maintained by Pat Pannuto at <https://github.com/ppannuto/python-titlecase> which in turn is based on [John Gruber's `titlecase.pl`][]


Rational
--------

I wanted a version of `titlecase.py` based on Go and wasn't initially finding it. After finishing the code I finally realized to search for `"titlecase.go"` and Google returned to me two matches. I felt a bit let down after my initial response of "That was a waste of time". But after reviewing the others code I found that my version was a bit more faithful to the process that was used in the Python version. That's not to say the other versions are inferior in any way. It is entirely possible they are better optimized for Go's coding style. In any case I'm keeping my version, even if the others are superior, mine is a bit different with it's small word list including _with_ and a few other touches.


ToDO
----

* [ ] Write unit tests
* [ ] Make sure ignore case is in effect where noted


[John Gruber's `titlecase.pl`]: http://daringfireball.net/2008/05/title_case

