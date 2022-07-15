import re
import Secrets



text = "At its heart, it consists of 13 hypothetical-scenarios, to which more will be added in the future. Each scenario contains a description of cyber incidents inspired by real-world examples, accompanied by detailed legal analysis. The aim of the analysis is to examine the applicability of international law to the scenarios and the"

text = "Eine der beliebtesten Kindersendungen feiert genau heute ihren 18. Geburtstag: „Wissen macht Ah!“ ging am 21. April 2001 an den Start. Bis heute ist das originelle Wissensmagazin erfolgreich und gefragt bei Alt und Jung. Von Anfang an mit dabei ist Ralph Caspers, der die Sendung zunächst mit Shary Reeves präsentierte. Seit 2018 ist Clarissa Corrêa da Silva an seiner Seite (zum Interview mit Clarissa Corrêa da Silva)."



line = "feiert genau|beliebteste Kindersendungen"
p = re.compile(line)

m = p.search(text)
if m:
    print("MATCH", m.group(0))
    print(m)


exit(1)

x = re.search("hypothetical scenarios", text)


x = re.search("consists.*hypothetical.?scenario", text)
#x = re.search("(?=.*hypothetical scenarios)(?=.*heart)", text)

#x = re.search("(artifical intelligence|AI|künstliche intelligenz)(?=.*bias)", text)

if (x is None):
    print("NOTHING FOUND")
else:
    print(x)
    print("FOUND", x.group())