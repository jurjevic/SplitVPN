#/bin/sh

doit () {
  ./make_icon.sh $SEQ.png
  mv iconunix.go iconunix_$SEQ.go
  sed -i '' 's/Data/Data_'$SEQ'/' iconunix_$SEQ.go
}

SEQ="0_0"
doit

SEQ="1_0"
doit

SEQ="2_0"
doit

SEQ="3_0"
doit

SEQ="0_1"
doit

SEQ="1_1"
doit

SEQ="2_1"
doit

SEQ="3_1"
doit

SEQ="0_2"
doit

SEQ="1_2"
doit

SEQ="2_2"
doit

SEQ="3_2"
doit

SEQ="0_3"
doit

SEQ="1_3"
doit

SEQ="2_3"
doit

SEQ="3_3"
doit
