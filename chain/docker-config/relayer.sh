rm -rf ~/.relayer

rly config init
rly chains add -f gaia.json
rly chains add -f bandchain.json
rly chains add -f bandconsumer.json

rly keys restore band-consumer relayer "clutch amazing good produce frequent release super evidence jungle voyage design clip title involve offer brain tobacco brown glide wire soft depend stand practice"
rly keys restore ibc-bandchain relayer "mix swift essence lawsuit plastic major social copper chicken aisle caution unfold leaf turtle prize remove gravity tourist gym parade number street twelve long"
rly keys restore band-cosmoshub relayer "clutch amazing good produce frequent release super evidence jungle voyage design clip title involve offer brain tobacco brown glide wire soft depend stand practice"

rly ch edit band-consumer key relayer
rly ch edit bandchain key relayer
rly ch edit band-cosmoshub key relayer

rly lite init band-consumer -f
rly lite init bandchain -f
rly lite init band-cosmoshub -f

rly pth gen band-consumer transfer band-cosmoshub transfer transfer
rly pth gen band-consumer consuming ibc-bandchain oracle oracle

rly tx link transfer
rly tx link oracle
rly st oracle

======================

rm -rf ~/.relayer

rly config init
rly chains add -f bandchain_dev.json
rly chains add -f consumer.json

rly keys restore band-consumer relayer "clutch amazing good produce frequent release super evidence jungle voyage design clip title involve offer brain tobacco brown glide wire soft depend stand practice"
rly keys restore bandchain relayer "mix swift essence lawsuit plastic major social copper chicken aisle caution unfold leaf turtle prize remove gravity tourist gym parade number street twelve long"

rly ch edit band-consumer key relayer
rly ch edit bandchain key relayer
rly pth gen  band-consumer consuming bandchain oracle tb
rly lite init band-consumer -f
rly lite init bandchain -f
rly tx link tb
rly st tb


====
rm -rf ~/.relayer

rly config init
rly chains add -f gaia.json
rly chains add -f consumer.json

rly keys restore band-consumer relayer "clutch amazing good produce frequent release super evidence jungle voyage design clip title involve offer brain tobacco brown glide wire soft depend stand practice"
rly keys restore band-cosmoshub relayer "clutch amazing good produce frequent release super evidence jungle voyage design clip title involve offer brain tobacco brown glide wire soft depend stand practice"

rly ch edit band-consumer key relayer
rly ch edit band-cosmoshub key relayer
rly pth gen  band-consumer transfer band-cosmoshub transfer tb
rly lite init band-consumer -f
rly lite init band-cosmoshub -f
rly tx link tb
rly st tb
