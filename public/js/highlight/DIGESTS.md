## Subresource Integrity

If you are loading Highlight.js via CDN you may wish to use [Subresource Integrity](https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity) to guarantee that you are using a legimitate build of the library.

To do this you simply need to add the `integrity` attribute for each JavaScript file you download via CDN. These digests are used by the browser to confirm the files downloaded have not been modified.

```html
<script
  src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"
  integrity="sha384-xBsHBR6BS/LSlO3cOyY2D/4KkmaHjlNn3NnXUMFFc14HLZD7vwVgS3+6U/WkHAra"></script>
<!-- including any other grammars you might need to load -->
<script
  src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/go.min.js"
  integrity="sha384-WmGkHEmwSI19EhTfO1nrSk3RziUQKRWg3vO0Ur3VYZjWvJRdRnX4/scQg+S2w1fI"></script>
```

The full list of digests for every file can be found below.

### Digests

```
sha384-+cjJSDOWSAbotsSXpQboVeVHgAguvM2DEPMB8jVhRAHSkx3tOhtw3kwIi6aLwrix /es/languages/accesslog.js
sha384-vx1I5PIkVSkAHbjcpmdlPQm/QmGrcmaqNFopgmjZIaCXdiSOxL84U6Vcccfmf+p8 /es/languages/accesslog.min.js
sha384-SXyBDWbbsvVyOn5TICWCahSdE/E+YqByrDZhk5477SMEixaEUJr8NzFO9Qjyhsin /es/languages/asciidoc.js
sha384-t9QnTos94FrBBrAC7s5wV5c/I48x4nw+3vsiNa9XwkyyQih5RbUUSqt2Gd4uSfkG /es/languages/asciidoc.min.js
sha384-u/MnWLPAdhUSLvSy9g3lRTFAwZSiveQjbKCTRw58BJh5n/RQfyyiwweCohM8LeSM /es/languages/awk.js
sha384-UjBtWlHRGRibzEDg1rIFn+Q1OGqWCyMh3RkUN9+hDW27sIB2j8mcGvxmkk8l8LVM /es/languages/awk.min.js
sha384-uBlc/xEFeDxZmBU7K/YWwi3ryXQrLQCAY2K1Dl3OD2DaAQBZZTt6Ew3aeDP20ix0 /es/languages/bash.js
sha384-4qer4rJCVxZjkwD4YaJfOnT2NOOt0qdjKYJM2076C+djiJ4lgrP1LVsB/MCpJSET /es/languages/bash.min.js
sha384-EwAFcMMnNYbOS/ABVtPlSLENfB27fQSxa/eUw9xWLejr7sWWPcICSoHJUxNVwhgI /es/languages/basic.js
sha384-3EDvjyClCK85gMCNUIHh/H9+K30vwYJjkGmjdUKpG6bGxahP+dSc+HRtXadT8eUD /es/languages/basic.min.js
sha384-0OZaeLK1yb5eP3nW4y0JP1fVharSrsuv/1mkI/6/8aRFm9laYIWIMXjCOqu+vRW5 /es/languages/c.js
sha384-G7WtwjMBNGrPly3ZsbPhfnT0v88hG3Ea9fQhrSWfzj7cNAsXsU7EKMQLyLGj7vTH /es/languages/c.min.js
sha384-hnJxdw/z4g01BcqtDA+k/+BkqmgiKrVrOi6NlogojekZCquQXhedMCO0wYSt2XQS /es/languages/capnproto.js
sha384-j1rHZcs0E8MTvcbjBDvgPJ2IWTlTwy90lyS2uAcDbDYEoMjL+OUyqaAPzN8Fe7YN /es/languages/capnproto.min.js
sha384-cEeSepkmM48Qe83yyTwFo/xnYlyD4ZO1tXaFbukowK2pPsea3XuCQSOwKnVmDrvC /es/languages/clojure.js
sha384-trRP8CnzW/SSG6kXpTUGPBhUS/tOwRQ992lsIs6/zZ9FueaoJeINicKwxLLVIZKV /es/languages/clojure.min.js
sha384-XT4oRe1NeFTZXGmTo1gp7cuec9Ap/64f/vt4I5GMRCj3pobm/sNDBZMhJsK/xlBZ /es/languages/cmake.js
sha384-O4LQvDBqImUEPey0HB0Cg3aMyXEf/P4KYm+7uDjqqeU0+ciFh6wEQd3jkRbLD7G7 /es/languages/cmake.min.js
sha384-Wjwd1YEG/PYlkLHTWIx+RlPK6XboMN3bEpveERJ8D8Z4RaNE02Ho19ZFrSRPGi0j /es/languages/cpp.js
sha384-Q4zTNH8WsDVdSZbiZtzWS1HmAUcvMSdTmth9Uqgfjmx7Qzw6B8E3lC9ieUbE/9u4 /es/languages/cpp.min.js
sha384-Sp8E1Lb9fNhbvqBiogM60TCgpAIwYBi8WbHhIHcXO0bR5Z+9LeYwpDa1gkjzU99W /es/languages/csharp.js
sha384-huUb4Ol37G1WrtGV7bn1UArXcJSjD4tQswMGzgpNZYAPxR74MHTqW79z1dXWMvhk /es/languages/csharp.min.js
sha384-DA5ii4oN8R2fsamNkHOanSjuN4v7j5RIuheQqnxMQ4cFnfekeuhwu4IdNXiCf+UU /es/languages/css.js
sha384-OBugjfIr093hFCxTRdVfKH8Oe3yiBrS58bhyYYTUQJVobk6SUEjD7pnV8BPwsr8a /es/languages/css.min.js
sha384-nJZn4t0S58K5MJ/2BOmgkTC8TtBWB76ykEKIIzb8sUxT3h7ouZH5Ucqnr9MJ6Bry /es/languages/diff.js
sha384-OIJrZgyhN1E6e29B9M3gFPtkR5JgXpls8h5D2bNzBxQ8Ncwxzfjedu3/0jZEe7S9 /es/languages/diff.min.js
sha384-hYnh6bxVnenrhg3HDafFjUPaO4yMWcNa6jjBCbcRJAzOW9RidrTLcmNgEyk1EkqC /es/languages/django.js
sha384-fYhQZafPWO54UUJ6/TWCLZVfptppoJrgZw8syHnvkUeAf7+lAiBkU3x8UQ7rTmvB /es/languages/django.min.js
sha384-CLo3x4rhnqUonsHs4zrr9gaMxc3gyLcDf2txgYmwyXDlHVe+FmQnYHguuDCbBQ5K /es/languages/elixir.js
sha384-B1unRJpfius4v1CDo8bNYlnKYL4h9uC7uhSE+fTC4WNSZwBHfHbD47N2bFBurjFF /es/languages/elixir.min.js
sha384-EzHQEbAzbwb1C9L6vi8Q2bMXMmOj4MKOUCn6OSa2NzRqDqVwGXfSVrylMO9zg3Oc /es/languages/elm.js
sha384-TpeiEEW/MM+zgu0X/O26T/xMSJtm4DAot0uQlZ73FP66AgdlqWO38C/GqpVwo58Y /es/languages/elm.min.js
sha384-s8ryKhOndFgWWigUu/RzK6ezBvWjCn+18SG4GMGLdTMSLeyY+AsN3fHP+etISCkH /es/languages/erb.js
sha384-4gcpw5sjrwU3Ojec/fllRQd9pKao49meXVCGF9QNr8q9shPlWffEZzHgh44NUpQ9 /es/languages/erb.min.js
sha384-9YnIRp+bSKwmD2dBIRe0N3NPQeMl2W+55IoWmQGhOqTWSMisYohRgghBd+qR48Hi /es/languages/erlang.js
sha384-qLpZFl0uFVtxGDLGp2q2g5PEgEgt78MggMr+BKUdDtRioP5Uw6RrYO7VzUheZ3ti /es/languages/erlang.min.js
sha384-qEBJBMpLpivuXkjdfOfq0lG7X75kJAs7lBKNa+lwHTB5p/rmQdnqpaiP9jJ8A0zy /es/languages/glsl.js
sha384-M7Ef1sJ4SsC67Z22eDV84pAdQRp9VnynFqNtD8xry35mYsHET3AmeJXScLce0iJr /es/languages/glsl.min.js
sha384-XcoPs3a4/IR69b/Dm+Q+u2pZ1mkg4OLQY3nregS20Zi93M17jhVxokvDNVBKwocg /es/languages/go.js
sha384-27jMAcMfx5pzlW2ntRUz1R22f43tLLOnYyDHss8iJBUi23iVzYrxbwQKY+LPU35B /es/languages/go.min.js
sha384-xG7D+neOE4ESa3no3TB3dmUhFAzqhWs44rleL/wq+h4+fnINKJvI2e/73Hcx67ig /es/languages/graphql.js
sha384-sYfgkTuyw8lkHmGdVV7Mw1/C2jkRl3ZKZrSceJsMaqZAG2RNed+pSOfjwMMEIBO7 /es/languages/graphql.min.js
sha384-6chagswO2Kdpzx/ovk65U07pdppN0fIgbg4AYmfufy3hIn+8GVZMvUf+dG+NVAn2 /es/languages/groovy.js
sha384-tgywtm20fLQUu0uIXBr5KDIjyHpfZuyM4zKiakHHNsi9zinDSqPIr/FVjREnEWmY /es/languages/groovy.min.js
sha384-NnhPwIPAAfzjFT+H35whxNsJEhjqnp16iuqrhXt+g1+eZhX1t3RDI6/wVZcwXQJe /es/languages/handlebars.js
sha384-RNsmJ7WM6P4KLAb4uTYFI84hIr0xBEtHjJaciwHoAsq00NcBk87WXUstUEDJfItU /es/languages/handlebars.min.js
sha384-FsHE2Bl98x5At/aNSZ5anAW5nAmSwAgkunnTWCeXa5pXLOt5dm8/6zp6g8RwEJWB /es/languages/haskell.js
sha384-HPIAyvou0DrFHYTekBAbUaQ3UbjeD/WEEJ31SYlgvnEQP54tky7yggE8xQqdfLk0 /es/languages/haskell.min.js
sha384-7m427rj9OYl8KKObOk+aD0QhXmqXy6hU/sHGS6OUwW9hRiiHvr+tLTbLiSJYbIyD /es/languages/http.js
sha384-GdKwMgTrO0R3SPXhm+DuaN3Qr7WkqYmmTClINAGPuV/BUPE9WUI/WnTEntp522Sl /es/languages/http.min.js
sha384-7hMYwsnoFLW0Wnjv4vRnZxuedW1tCz7/pydl1b9w3fg7B+rldToCjqzXwCRHUE7C /es/languages/ini.js
sha384-LDu6uT3diI3jBUJtdoGFa787cYlVrR+aqFfmW+kW+TImOkpVe+P5LBdDzydIHo9Z /es/languages/ini.min.js
sha384-WFJdA9Hz+G9NQx5vPba/tcGyIibm57UkKVY32wNB/94iT2FmPma5W7gY8p2l6qps /es/languages/java.js
sha384-coaxfgI2lKuDqSxfMlfyPq5WM0THaLGyATZHzaFMrWdIbPcLdduuItTe6AmT/m33 /es/languages/java.min.js
sha384-WCznKe2n87QvV/L1MlXN+S8R6NPUQGU34+AqogMuWGZJswSD6rt3Mgih+xuKlDgm /es/languages/javascript.js
sha384-eGsBtetyKPDKaLiTnxTzhSzTFM6A/yjHBQIj4rAMVaLPKW5tJb8U6XLr/AikCPd+ /es/languages/javascript.min.js
sha384-12GbYFzdyZCSmfBTmO2CR/qE89K5uE1PEuJ3QUwXH0Q9u+uoLNigjH9dG7LAxxiI /es/languages/json.js
sha384-+DT7AtubDhVDciRc6CgjJJRsCt0L8NC3Dh8n9Tj2tZWU8rWxDIj1ViubmUDn8OCY /es/languages/json.min.js
sha384-T5TWWx2SVBqE/AJVqpKp6D8+8rpEcX+Usy65BcRR6SM8QEh62nMtxDsowBlhx1Ay /es/languages/julia.js
sha384-NqyA97ZywXJCu5WG4NiDyJRAYm5L2aGPPTeGnRSfkEtK8Lch/likdativWWAbLUs /es/languages/julia.min.js
sha384-tsX5LI0gdW2Xk9ZMsj7B2vRchm3jC0zoc1r99Z1377aXFJJXimRtRYQprUJpSuu0 /es/languages/kotlin.js
sha384-mfmdbPhLobPr5OJzSXlWHDQDymRYyzedurWjJBvKVhlQGE+Jz/pN3D9lPEBIkDK2 /es/languages/kotlin.min.js
sha384-JZ1yKn4lHbd05rOuRoUSoCM323D+X3s7vLYhESynzVpdWAFu5JNp6chDl5N9uc1c /es/languages/latex.js
sha384-dDJ50ltJ9VlfSOzSORI0KSNPcRIpLy1J0BmtdQGXBEY+7f57DAx754HVffzmWqC2 /es/languages/latex.min.js
sha384-rTufDmtZDbCaXsmgV+ApffWh+8kRBEMiQoB6BRsoKvf402NBZGR4xC2U8VrxHZ72 /es/languages/lisp.js
sha384-yR0v1Wvtd+8Y1TA+Tp5e+PEYgonMOQpp92Ox651UghHCxvLkEBAuEBWWCZqm9wGI /es/languages/lisp.min.js
sha384-OkDbf1Slbqz4CwuUJVZyhq9SKn5vYk5ADIxIBS3iV5D0ZY7T3g3BOGZgd9u0kyZH /es/languages/lua.js
sha384-wvrGnUzHwJzO9dWQFF2DxrFjkSPw5gmc0iOQYmJUzeZ4tqa+15VEFRCH4GI2DNF3 /es/languages/lua.min.js
sha384-UgQTewauLJ4EgpADCJ99JfEtiPvw+fyaSrY1gtCVBviDNG7yCH5U7qutYptSfYk+ /es/languages/makefile.js
sha384-aUXBqKsjOzPD/W+hccF21KKWmWts/CrY/lWGJU+dAcsKtuh1/XtyDnzfLmqy/fV1 /es/languages/makefile.min.js
sha384-bWwkdmOCj83zZ8/m+oPD9goRMhrPCb25ZA6aTyg7vcsq9IpuyED38kQSw1Na4UTZ /es/languages/markdown.js
sha384-SqGSUq0DMQ0OUQnQnTuVDCJyhANd/MFNj+0PF67S+VXgHpR8A4tPsf/3GoSFRmrx /es/languages/markdown.min.js
sha384-aE6FFuNOdWKOdWyCTVnCnZH4NBESWCBoti6+Jn5mq+Ss7DMPzJKHz7W3VkYTFhWz /es/languages/mathematica.js
sha384-ZOli511kLeGbbzCcuVX7mH7u8A+Vv5xoT3HaBPcSW8DwrsJkjCA7Ln+Wluj9Xm17 /es/languages/mathematica.min.js
sha384-KW5m+lJxeSte36eSdXdF6RWKt1rP3PKVAblsiDORE441O/7ONXso3W4FMiQM7DVo /es/languages/matlab.js
sha384-gU9Mv6FEG6RVi/d3cgP9RYHSdGMkLi/gYQCn7smdbY+m1qDCN2paZoONWRk0tBP4 /es/languages/matlab.min.js
sha384-1ToS9OQQb3UX17cXpLVMXMpe+9vg2bz4y8u65oNIsnfzRj4VzPl3uspOeoATmH1K /es/languages/nim.js
sha384-trp6hYonRewpJnUz13JwZm+20STvnBrJhyTEN8J59do8eREiQHY2UMDfbBR/L4EO /es/languages/nim.min.js
sha384-alo2ba4JSULG3Q6xAv7jIoMsn5q4oe1IEqLojzvRqMRekael5+kttHkoyYklPz6n /es/languages/ocaml.js
sha384-sm+3C97241bDJHul+IONVgYeQJcwtSN2wo1kUa1K0mNWY6MPMJnrx1giFrUr2MNY /es/languages/ocaml.min.js
sha384-B+EjdRTJfuRQL8OfQPSh3FpJXyftaETF0bCfUNAm34365HyQ0yM6cjh+8qYsSEE2 /es/languages/perl.js
sha384-kW4vu/krgDKKOlEd71XhqvOJNxZuQ9t5Rwutcxm3/zfujAXe08IcLMX5Al0NGE72 /es/languages/perl.min.js
sha384-JBkI+6623OoC1qCgG+MY/Ta0qRYSzTDH4NGMA+7U8RNOjkh7geFvYpRidvdHs3zT /es/languages/php.js
sha384-6Znh5S5q9F0hGNMtdaihFL3RJeMqL1a/UecjgzaBrhoxe4sjd5aiEOgokt7blhSW /es/languages/php.min.js
sha384-TTDGPCrk8Dg2oW6NghGM5WJQPbi34BdYJj6yfsDiGXlM5os/SebXT6KzATp19rzo /es/languages/plaintext.js
sha384-XXx7wj9KPm08AyGoGzzFKZP2S7S+S5MbKMPnQcWUyhJ3EjHvLuctK/O1ioJnG2ef /es/languages/plaintext.min.js
sha384-PpU3S+yZJ1Gj2R7L/OgTgkKELQOf7F1VcbaAEJwSFSM9lw2ON2wxq1FvUcUlm9ne /es/languages/powershell.js
sha384-pBczxETXXX/Ne92jpBviT47DPiv27FLNtvs0aMZH3W7rYXdwgf1gWg3pH4NmNC2V /es/languages/powershell.min.js
sha384-m+uqtvx6pG1ygzVw5qHomirrR8wyyHi8RKEOthLbwHmOuHU2uXDask8IcbVVHWdC /es/languages/processing.js
sha384-t0XeJ1PaBwN1QukCYAP/EPBK7SXIMRxw2R2nm/P15RGqdxjFIFmW9+7i/Ilm9Xqa /es/languages/processing.min.js
sha384-Qb0lgsDvkX4k+zBLrQYk4QRuYU8yxfgXrhCQaha+21MnncDVhS7gnQIU/CI39I9n /es/languages/protobuf.js
sha384-93TO3rpvKyXqikneOt71cSHkxWQWPKaa6n00Me0j5Fzvjb7j+EjURqbZHyqx8m8T /es/languages/protobuf.min.js
sha384-+5oyk7Ed3OlvEWGj8xracq/6e52BScKUN4kxcreNwB7kfRTVsAMs/aAJM58dzIFN /es/languages/python.js
sha384-ND/UH2UkaeWiej5v/oJspfKDz9BGUaVpoDcz4cof0jaiv/mCigjvy7RQ7e+3S6bg /es/languages/python.min.js
sha384-GLWbnkTbelVKeKv8VxFhtauUDfV8uqxjBPEEZ3WYks9zmbd1rFCWmdzGucZn0+ZQ /es/languages/r.js
sha384-T5hv8+GOm2ni1B0cuZtt+8T0cIMH1CTycc9OQQH0GLcIQhmeBosyhcmcIv4r/1ag /es/languages/r.min.js
sha384-vs4jXytP3N3d5zRX15Fqc4u/kDI5jDr3PNYIvQ0mmmoU+o/yxJfu/+VYme6qIOa2 /es/languages/ruby.js
sha384-O4a+vELXT191NhKLE3TR5WQwDmQZ2izAhb2zETRxcPSXr6CJruqJ4a+GJpDlaqCF /es/languages/ruby.min.js
sha384-I4aH0szMeaCbcs8R7dhxA3p7ZBL/HFxnD5Gbz6l52kIrd/igSSFi/9sJCykNuL52 /es/languages/rust.js
sha384-1vvSh2x0WCtPLbkTMqNuf8JSZw8N5bSo9oONZ9vqU8NOBHPIuKt+kFdC8G5nA+P1 /es/languages/rust.min.js
sha384-qoJzsMEj+/mLqO8ijrly9EdS7NefE7hIiEwpLw4vop1thdpo3DbXFeQc8hYnrtFR /es/languages/scheme.js
sha384-lEVNmnotVzJyk7ga6toEDtsBpgn9WoqwbsxeCyDEsmdGqyDBNeWI1kFiqtJ8lJF8 /es/languages/scheme.min.js
sha384-lRhSX2XDrY27NzrAS1t4YaeRtwjsY41kFBbIEYltkmnsfSE7lbBJMQVds2u/MqTT /es/languages/scss.js
sha384-RDUehV4j9Do6iGkYq9Gjn3aUxh6x+NFER1sHpLUXsNoCFjah8Ysrlad8ukLbIr4J /es/languages/scss.min.js
sha384-dzLjhd9nNJH62idgKI1vZEKHRBtZXSgwWQdPR+emG7tfAN4BW2g+A5Xs2315Uxii /es/languages/shell.js
sha384-RKUoelG22/D7BV/bNpoGLNzdTgWRf/ACQX7y4BGyIzK6E+xUoXtm68WNQW2tSW8X /es/languages/shell.min.js
sha384-rBAFhyrcRcMNbVJ9g4k5y3eQDkjGdgoOb3oTWTbHgwyUgUNv3CK9wLsGy/d+52oa /es/languages/sql.js
sha384-8G3qMPeOeXVKZ0wGzMQHgMVQWixLw3EXFAcU+IFNLRe0WoZB5St3L3ZLTK62Nzns /es/languages/sql.min.js
sha384-LpDkuXFg1D+54cHhYqk9r9E4vKH0CGAnyBqiq5A6SnvPEIMTkMH8IN8i14JpJNhs /es/languages/swift.js
sha384-lJj3aAxzUpdk8StXN5j3OP20/Loadv+t8jYdMBYVqCaxtLHTpBUalFDsTPkC9Fov /es/languages/swift.min.js
sha384-gpOkV8Kb7AIZfX8p/6dylHJqwwQF9XPYIIKAmDAdzKhTZWFYwod2YVYx+lHSatyk /es/languages/tcl.js
sha384-fdDOAbgCCSuCQZYLNBFr1Iq9P8nOC+SR8iBVSCfMbSAHuZjNmqZjtW7etWgoOHZN /es/languages/tcl.min.js
sha384-CS3qiWid6Sod3yAiQwgPzy2ZerR00u/cwhnMxQrETuI74o006r1p5qj1U9Gdo4uD /es/languages/typescript.js
sha384-HHgazax8ozQ9RDWlJQyaFzaTk/OgTFRVHH+lcaYInkE8wVu5fnpkqkB3KUdzKcvE /es/languages/typescript.min.js
sha384-C5Ib6vLvFc6BpxyyN5165sVTrN05Avl4ADMSVNVSmB61oPr42vOavsZlxOknqqWz /es/languages/wasm.js
sha384-FlAOV+6BVjUyajzW4twmqUvHxDYallf5jw5FewMZmsBHBuiJ/CJj1NZhNT1fHuzT /es/languages/wasm.min.js
sha384-OFoR8IZ+CFwcY8plx8HSDZNoCwLxc701CwdNGfoIEhSgwAbwhvInaxnEi3HYTt2Q /es/languages/xml.js
sha384-yFd3InBtG6WtAVgIl6iIdFKis8HmMC7GbbronB4lHJq3OLef3U8K9puak6MuVZqx /es/languages/xml.min.js
sha384-MX3xn8TktkjONV3aWF6Qn6WZyq2Lh/98p0v7D0qEoJ4WLdYjoAyXF/L/80q3qaEc /es/languages/yaml.js
sha384-4IiaMbQ0LBiRJYBGoAXsN+dV9qu/cGLES6IuVowdeCu/FAMY5/MQfD/bHXoL2YBb /es/languages/yaml.min.js
sha384-wWNyZqG6+Tb1rRskCP3qZjTwRbb9AXwkj1emJa1QFm/wqE5WZf93bjR9IAgcxXD7 /languages/accesslog.js
sha384-BF61I4vqlKI+Z6dQUyVuqQEE9RSYZWptePUuYrYLfxlYJT6r9qNemzfyAqJIXHRJ /languages/accesslog.min.js
sha384-HPT+hE+CRAgt+UgsOWKsXyconKSjRIbAxiWxwmmQ3L3ijW6pi0tvaMpnT0vNORvu /languages/asciidoc.js
sha384-sf6kQzw438VELcTHWkS4YdEyIzCPQbGZLke4v5rtAIBrnGsuvCILPWGeIcDIVLmX /languages/asciidoc.min.js
sha384-xfaX9UdEVoIlYoDcsxk5OO1ddyQ0fl/BmWflf6SctNce+UzNDEimoWGmLbk0FTPg /languages/awk.js
sha384-+kQeWvlXeWctdKQQ8iGrG1JdmiKRJHAqofS2vKNnJ24H25HqMhb9kIQC3NYrsoch /languages/awk.min.js
sha384-qbbaBGYYg7PdopdWOGj8KdkBosUDY5PAe3aTMJKTqWcriPBJJzCVu5BlwNEwqr9U /languages/bash.js
sha384-ByZsYVIHcE8sB12cYY+NUpM80NAWHoBs5SL8VVocIvqVLdXf1hmXNSBn/H9leT4c /languages/bash.min.js
sha384-KNVj6C9bNyJrl0feAi64iBCZkXxoTQxDbAJqGlVos5lGlSTOHd+zE1xQgkPFoH4t /languages/basic.js
sha384-KfOKKkbUp6gWmUdMPNCj5Hqqmf9Z03BTs0GnM3hyrYlNBYXSt4O527wBkKrgxiIl /languages/basic.min.js
sha384-VZxKf0mjKYDwZIgrW+InqDfJ0YwYUFEMw/4YmpV1oKtVXFVmVq/Ga1vgq6zKTUqS /languages/c.js
sha384-mWenE7UXKnmYRcJ3mh+Os3iZ43BmFf9x3AZMM6gi/2sT6vxouPWspbTdREwWDO3w /languages/c.min.js
sha384-9HvxzusD48ytHDhXr0xe+uagRLJD+Er1RFSnE7tH/ESuVNvtQHivWTBSk7cIYIeD /languages/capnproto.js
sha384-X9ejkQYTr8UwtUbCNj5TNFyggrCb1iV8lydyf6Dz1OyGJRuVopefQTJOK1KevmS2 /languages/capnproto.min.js
sha384-k+XDUAVhlUtu9Gv3AkKvB9gL0WtZTonn6fFtlSdKn1+7U6yNLgRLHg84pXr4hmp/ /languages/clojure.js
sha384-01/jwK07TwTjPcx55RJ1V03zqCdkEK8QBmgLSbcgcyl1uZ1AwjuSmufE481pF0nY /languages/clojure.min.js
sha384-/DE6n+XheQ/kdrhU7A+cWj1GSRUM44zfgDDcWPvIeGGr02ckG45tZIoLbVgx4X1h /languages/cmake.js
sha384-iqJTnKrAPlpoueHQticEO1SgUDmD+Lv1EGplrE41YfNgnnVCt6pLeVjutOnr/Ddi /languages/cmake.min.js
sha384-J4Ge+xXjXgzbK2FP+OyzIGHLfKU/RR0+cH4JJCaczeLETtVIvJdvqfikWlDuQ66e /languages/cpp.js
sha384-LMyrRAiUz6we2SGvYrwDd4TJoJZ+m/5c+4n4E64KVkfWFcZdlrs4Wabr0crMesyy /languages/cpp.min.js
sha384-8sbRwiU8Ar2M7+w//1u7YiI1e7KsmB4k3QbW/m1IW5FVH51HiOpR7g5QGE3RqTNi /languages/csharp.js
sha384-wWP4JQEhRVshehTP7lUMDn3yhDI2+398vN2QW5LBt1xIpK0Gfu4dPviO8tP9XRo5 /languages/csharp.min.js
sha384-r9czyL17/ovexTOK33dRiTbHrtaMDzpUXW4iRpetdu1OhhckHXiFzpgZyni2t1PM /languages/css.js
sha384-HpHXnyEqHVbcY+nua3h7/ajfIrakWJxA3fmIZ9X9kbY45N6V+DPfMtfnLBeYEdCx /languages/css.min.js
sha384-d6WNn2kPYgdezhgRmGehSkGVKwbg+ImJrGTA43yYg8VsBqYW5TakSDi/jbVjzqfD /languages/diff.js
sha384-yE0kUv4eZRBzwnxh0XmeqImGvOLywSpOZq2m0IH40+vBtkjS716NYEZyoacZVEra /languages/diff.min.js
sha384-dFC7/UlAe1aH832WvFmt6fwQFIha+bFdz4Jw/Stp0m65X0P0zgiyaSYVKpRyPCOo /languages/django.js
sha384-AqTOESQu37Lj9i0LQjA1B5Ju2XOJzwB7RR1NYcpo+JgnUF+UTdQi8km+UzU8uYBZ /languages/django.min.js
sha384-Bd2ywvYdN43SxOMuK6TJ8MEi/51znKu0cait9/ssPMJi6p9Q74aDle3iJ7V4bX2o /languages/elixir.js
sha384-qMYXRUi/Q9Ukntag3SOX3rOfkgLmmjlGA84X26QeOGixvWyfSGEiKwaw5ZKoSyBH /languages/elixir.min.js
sha384-r7qcABiwqGhFWIM6ow6S2NwC6KkQtyvyl+UPsL5B/7GHqqqkEbfJOJEAsqzlSY5I /languages/elm.js
sha384-IKAy8uG+2GJ7OKT8sMCyI+irjkHsROGLhUVfVtbTXwtjibq81wfzBddeahMGbprC /languages/elm.min.js
sha384-OwWPF1/FA3TnWxXB8Qg1O/DBXv17fMJNHsXv4ocskWfPvnQ+ztmBBiQrYOYd83GC /languages/erb.js
sha384-HH6ifvxwyuXXyJkqfZjOkyD/ZQIXMe3zZkjUKZKwCDcpbiiq+jLNyBEclXTBwK69 /languages/erb.min.js
sha384-l/6ahuMmod04uph7gyXcv08JuDvo4gqXigxibqkcLetdKzted06RdzuxDk5MO0i2 /languages/erlang.js
sha384-hsThqIkkHew8Pfgbao+v2vLaYUKOUiN8CiPSkJRrmxTe4LcgVjZLC3Rz4z2zMT+X /languages/erlang.min.js
sha384-ZJ8K2T7Eu40AMSL5L2melYeL6dBE/QobNlRdRg04totelNOYMs3d2igwh2KXNXNC /languages/glsl.js
sha384-PkNQzkPieB8PCovDH3VAfEBO1RvFP6kPVFxLkO2soZ3MoAaf/TeKffj2vMLGxJlT /languages/glsl.min.js
sha384-lDCjdnxlW/GRZYzy4Zqkj833wJD7Hc2FP927RAtySEdrShMiUSXsWuFy1IC02qxU /languages/go.js
sha384-WmGkHEmwSI19EhTfO1nrSk3RziUQKRWg3vO0Ur3VYZjWvJRdRnX4/scQg+S2w1fI /languages/go.min.js
sha384-+9Rtg2Bz4pdOgkMjD/Y4IbvMzkuZqZyZwhBIbRD7254eY5Zm2Sf3UFI6Bt5tpyT5 /languages/graphql.js
sha384-yqNr+JW52pR9ZNw8CIHWgGMIrqhKhuOBDSlQ1ulf4Gt6wqi+VUMGHP5j6ReSSRY5 /languages/graphql.min.js
sha384-ow3d32gzhvqfj5wlYXK68qFqHd99nym/Rmfwu/JegRYEk9tpbiEV38q+STNcNMKJ /languages/groovy.js
sha384-N2BpqybE0W3/jhUPUCcAr9I33RlAVfy/E+lEPCSY+NLvXJdrOUWK5+/ocIoPmA64 /languages/groovy.min.js
sha384-pk3Yez/JLDM2LWCA5AUcpHI3zOTzwIMp21HT9Ugn6Cw3ZKgv+3UDmcplG2wJpHBo /languages/handlebars.js
sha384-PqtEiFXspmwbGzxcwaN07nChfbxIAJHK4bTNOnytspmN9Q+LkL3AjxV5SgUI0WJQ /languages/handlebars.min.js
sha384-SCOieukbjumPuUXBq0mCKyJNv/n4utN1U2QlTE8jfF05lCKgWJ4SV+SDp1dHgQzd /languages/haskell.js
sha384-EnLqy/LmA2yShypXJIP5cvUz0S3jJ5l3PnZcMRE6JQdpeNOttU97yyQAsdTusDbl /languages/haskell.min.js
sha384-g/lhU6FXH73RQL4eFwkPk4CuudGpbHg6cyVZCRpXft3EKCrcQTToBVEDRStjYWQb /languages/http.js
sha384-wt2eEhJoUmjz8wmTq0k9WvI19Pumi9/h7B87i0wc7QMYwnCJ0dcuyfcYo/ui1M6t /languages/http.min.js
sha384-JSUMPR+WT0h/7NlqXi1Al9bVlNT31AeZNpAHttuzu+r02AmxePeqvsZkKqYZf06n /languages/ini.js
sha384-QrHbXsWtJOiJjnLPKgutUfoIrj34yz0+JKPw4CFIDImvaTDQ/wxYyEz/zB3639vM /languages/ini.min.js
sha384-pYIeBYeCE96U9EkPcT4uJjNWyrB1BKB41JIadYJbvmGa5KacaoXtSQOUpBfeyWQX /languages/java.js
sha384-uUg+ux8epe42611RSvEkMX2gvEkMdw+l6xG5Z/aQriABp38RLyF9MjDZtlTlMuQY /languages/java.min.js
sha384-vJxw3XlwaqOQr8IlRPVIBO6DMML5W978fR21/GRI5PAF7yYi2WstLYNG1lXk6j9u /languages/javascript.js
sha384-44q2s9jxk8W5N9gAB0yn7UYLi9E2oVw8eHyaTZLkDS3WuZM/AttkAiVj6JoZuGS4 /languages/javascript.min.js
sha384-dq9sY7BcOdU/6YaN+YmFuWFG8MY2WYJG2w3RlDRfaVvjdHchE07Ss7ILfcZ56nUM /languages/json.js
sha384-RbRhXcXx5VHUdUaC5R0oV+XBXA5GhkaVCUzK8xN19K3FmtWSHyGVgulK92XnhBsI /languages/json.min.js
sha384-ejh0O6l/Lf9xflCigwVR6FqqJAWhVWB+M7kjlcNSQvtso+e2QBqicY9b56ih9a2X /languages/julia.js
sha384-ZqjopgKriSJBeG1uYhjsw3GyqKRlsBIGaR55EzUgK9wOsFdbB67p+I0Lu9qqDf2j /languages/julia.min.js
sha384-UjoANPpyYDKhhXCD7M94TUBIvEg+Yey//mTGwAKIPkAcmpT4QVk8D47YYWumt/ZR /languages/kotlin.js
sha384-vfngjS9mwSs6HkzR9wU3mDDip7sa8TLKAxsuQ9+ncUHU25slHzHOdx/0FWxvbg4I /languages/kotlin.min.js
sha384-0DF4M/uGwA0TYtI4/4/Nxlz/Z6TObHvf/Uhz39xsvP8TkITbNVOVhbxyOXgyUn2V /languages/latex.js
sha384-0ZgcpgIDZoDPMdyG1FHnl/MByMKt+W52UltboQpuA52+f4pniKOMFXfIfe4K4Byc /languages/latex.min.js
sha384-BlBnQIfBY+Ul2GpGcD6YIGJM0Ro0rrsR6dmWzQ4INqxG74whxz/Peg3du4VkiOXY /languages/lisp.js
sha384-+LHHMbAXOUlvquvrQZ9LW4KeR2nwcsh/lpp7xrWu7KuaDSGgAYBIdm8qCw97I1tq /languages/lisp.min.js
sha384-/Ml+gzp/rkQcZkCwBqpVjCj028T6aTnOF/LCRJH8LBE5xcPcTkntQwJ5KGMMNLI3 /languages/lua.js
sha384-T04Yu4dcDCykCMf4EbZ62u3nURYEVkpphRGFhF/cMu+NrtDqoRHgbVOZz7hHdcaO /languages/lua.min.js
sha384-LmfE+sO0d5qZL2Ka0DIrgJ/5U1plo4uFFAmgjcMxrQO+RkeWVYWuaphHAdrY9g7V /languages/makefile.js
sha384-NIrob3StFQyD/nlOsXVCeRsJ0N2SvFEDjFtYS393wbD3CY1eT+2kwT4RL7tpMMhs /languages/makefile.min.js
sha384-U+zIQPoVdPCO0o4poik2hYNbHtNm+L5OojDTulgIeEZTNz+LooLAm72d66mNjwKD /languages/markdown.js
sha384-mCUujHHbWJEjcupTTfWOk9YR3YCYNHaA578+TTXUd4LPi7fGNuMQbysbl1pmcIGd /languages/markdown.min.js
sha384-pwYx5UHO9OuH2fKa087IpKIBIlWTlzTGYo/nUJ9C2vGcOy19NI3zxCU/LV4QzL+N /languages/mathematica.js
sha384-Y0PVkUaoJo8lSq1gzhxTxCx58Dd5w7lB/0RKYuIK7hmiVyndY9vplBBava32RP+m /languages/mathematica.min.js
sha384-NKxL9la40gyvGCP5oLrJcIJDWXXv2FJc9aatPho9P/PAH1IcOBIcklePMhGDlI6x /languages/matlab.js
sha384-KveG8st4Ls4iaD1XzpsBzVc7g2K6kxnbxlZOy9cH6Knla0ZH9jsloP8nOOs6WYMP /languages/matlab.min.js
sha384-6liNm1UFDJnx7SCXza24sp+BKJBOaFyEhnZi5JlPK705N7WYpT7LgSDN+T8SgKeH /languages/nim.js
sha384-aDGSumJzWocIW8uUlSvYbRE1BwKU319shHv1p8V6DsT5RNsImzghku0TxI+MrU1E /languages/nim.min.js
sha384-YcK6Nd3iy6OR5ha7IpqDnBQpa/urx4p84U362siGzBuUqjqjEgjq/a2TX1hkS/gL /languages/ocaml.js
sha384-WhlOQ8tkNAn2N+yTa7njealQ8ue+22g0Zw/iR6A9QtSz0XQZzBebXuAvINaGCpdY /languages/ocaml.min.js
sha384-YKOtRjfzHo4SHG+iI5acI0B8GrsAlqbQpZzAKKR/oABDfo0Qw6fdjwWtHY1He+ih /languages/perl.js
sha384-J8kCjl82Po/VNpumPop593F3EckTvJlnUS2ugGW46Zjts8BPGl6iNkXjcenKcwCx /languages/perl.min.js
sha384-S1JDGPScVg8ikNKLZc4CSP0ZxLiJ7bOJMzTLfOzQiCxR6wPqAa+YtauHJXQpc2GV /languages/php.js
sha384-c1RNlWYUPEU/QhgCUumvQSdSFaq+yFhv3VfGTG/OTh8oirAi/Jnp6XbnqOLePgjg /languages/php.min.js
sha384-IHapUcPkNR+7JNsR+qYSVYGCE3Dpzo2//VYWtmGYrw3eQG1RItQ7HYq6aK1Jo/6L /languages/plaintext.js
sha384-ofjxHpechXkaeQipogSyH62tQNaqAwQiuUHOVi4BGFsX5/KectIoxz16f8J/P5U0 /languages/plaintext.min.js
sha384-BW3J/7iREOC3EuqJwzGkDJb6FY9bVvkt70ipieV892fYAp8uReDMiO0/lrUPmkzL /languages/powershell.js
sha384-LWJZQx0dgGhEK7snfNYrQs5K+QKD1sOmE02sOQCz4br9UmqSJDvPLoUVFaUyFnjq /languages/powershell.min.js
sha384-xtY2My03c4KTTlpV3P5oCctL4S6Jwqbm2sqyVGX5IA2n13I/D+uVXX4ZRu0yqZJA /languages/processing.js
sha384-Hu1K5IR3tWjqNzVwTlOnFre8rB9qJ4qjSVzWmBVKbaTtCSYxFRPZrlRHb0chbeZK /languages/processing.min.js
sha384-v2nhli5f4LiO6PnSA6zJgJ/YqCBEYjpUyadALhzMSH8h9cN/iFJEOTFOtXw9jzda /languages/protobuf.js
sha384-1ALwZ77/iqBsHTR04APuiMf4wWnFgY4b8VUXQHCjPfywk7S9K76uwvs1xLZUTfY5 /languages/protobuf.min.js
sha384-zdZio5RcGiKQJCpe/1IXujPle3bIY8sbmvCabSU5G5GzWAzZtoRZfg9QAQXCL08q /languages/python.js
sha384-IP4vv4Aoh9Lyg8QyzVkAmn2JGoDCpgVHzVSrD3Z+rVyn7+s4wx4pRjv+go3TEwfj /languages/python.min.js
sha384-Xa3MEXJyV9PD/wnbdHTdqvipEl9QqodA35i6Nmw/6wqrvQ585Q8EDfz7f2UCuYpU /languages/r.js
sha384-+z7INDMj7HMs67TuLfoyKa6TUHdTweahANtJ1Spg6uCa8/f1jFUf238jj53TDxmP /languages/r.min.js
sha384-AE7f30oAuQtqmFee9hPd2uoo44ZkyUY2wQL7DjjENhh2wrvS6q4mpxGtAxeCxiRQ /languages/ruby.js
sha384-pFZpTUpdH4YEXSenc9hfKZ4uCv2IQoJQCIlIHpA0fM2cvTVH8LuzQMNcGSRGeJG0 /languages/ruby.min.js
sha384-CA6FQ5i08WYjgGIhQBrXKmcJg42apGjTP9b5WqttVw3cYEtXwHHGo+XJLYS7u7F2 /languages/rust.js
sha384-ZQJ5PCEftpFqCZkLDs96CSDGddxBultwqTdlxjnJ5h2doMAQv0n1x66w7T/JQEyy /languages/rust.min.js
sha384-7fRWpXtyFWskPWp6k73jlKTxGBn4B7kRRMq8eFzjqnmTQDJdQmvVTGQjNfHQeJbu /languages/scheme.js
sha384-wDAOgnGx0t7X9tTKddfqp7kcfAYVDaQDRdfDD+47IukNBVuVGFldV80H4qDTHPOR /languages/scheme.min.js
sha384-fwYddFsITuK2bPhi9RuIzwi4PTULEXgtEJsQzTdx97vOS/GHfrk+aNSLxEHgzQa5 /languages/scss.js
sha384-6u+QpCDqQidb5pcO+yBqy0xLJ30x30VlrFvXm8J84LMwGIw9q3U4u+Z9vFXlhB5x /languages/scss.min.js
sha384-KW3ZDReTAemYUfVHvH1MNQ/v6agCYYdMGdMteP/yVV+NetIJeDMx0ruUMTbr/SD3 /languages/shell.js
sha384-PDEMny3uYpCrSKwqVGLAwqCsc9u+9hYXauxVPqO6RXmlxSNKLoRByEBigU5lahmK /languages/shell.min.js
sha384-Dy7I/j0yJlyliWiNrkNqXfxDrbN65q40s3JColgTYZQ7QJa7lcmK0WUL3i00/T51 /languages/sql.js
sha384-8q00eP+tyV9451aJYD5ML3ftuHKsGnDcezp7EXMEclDg1fZVSoj8O+3VyJTkXmWp /languages/sql.min.js
sha384-Q1hSwaZr828HvfGGkcH9/K6Cg93VFYOWEZ8cpL7bUdzvzrSurur5RjoZjn46PokC /languages/swift.js
sha384-CYmrQ9dmDVxuVoM185jHQsjhiLlG/kQfabzRdOYsfUV2AQvpjQNrd2zVCpCC7N5j /languages/swift.min.js
sha384-Sp0SaiUveBdlpsayV7D5zyAQqg8sE7eC3j1fYZ6rLr0k0wdpmU16EVIyTdQuOLhk /languages/tcl.js
sha384-88ARtX8chYVIZZ4A0cFUpXilu4qORXSqH3yTVPOSoOSluZvBRZ6TAzO22TvW9Kg9 /languages/tcl.min.js
sha384-yZXtQC/OmWoPykosK7vE1nCvV4E/six6+apjNau4JwBkejkea5nP7VBEJJkGnvoF /languages/typescript.js
sha384-ORwtVEfrCZ0gzGacgmfv1wOtxcPIaVfHKwq8dKQjObRwx3qpKjsSg1ldTu1PEgXd /languages/typescript.min.js
sha384-ry6bgnkpMB819ITiAptb7x+iphp9TkAJTFl2NZQfF7JCOxZE+DcJvsAUvP2t4Ccn /languages/wasm.js
sha384-fWMR10q8HCFcv4u/2S/P5JHw477fD1akWW/mGsZNig4vAOk4135GEWJft9llui8U /languages/wasm.min.js
sha384-+PuZYFfVX2UQZU2yKt/FsJUZNUPzZWxW7auXltsaecr1xLvzBYF3c5gYoyOs1++x /languages/xml.js
sha384-jgkY4GMNWfQcLLIoP1vg3FWXflDrRhcSXGBW6ONIWC2SOIv5H1Pa57sXs+aomCuZ /languages/xml.min.js
sha384-tB5cwwsX4Ddp7P4d+ZInDb3nt4ihEEglHXoQ18eVLlT7soEn7bfGfABWKIn1l+H2 /languages/yaml.js
sha384-WC56y8OaFPt5Kj2HX6JAumxUYEjQmBDcSTJy2pn/N8g7dg1hKjeNVrJYoxlpeVmz /languages/yaml.min.js
sha384-ydNcuUkuj5IOH+8AyGqMMhNqtl/HrQwFvic2LaKzZwZFtqVGPVKTttsAgXL2c7fW /highlight.js
sha384-eazGvd5GTMy7Pn6VfutumcYjFyjPCo4BuWhiOFB6IeoqKRBJC2CZknz879aULnMD /highlight.min.js
```

