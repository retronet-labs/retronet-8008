# Testdata

Questa cartella conterra' ROM e dati di test generati o curati per verificare
l'emulatore Intel 8008.

Per ora non contiene ROM storiche. La CLI puo' caricare binari raw e ROM locali
tramite profili macchina, ad esempio con `-profile intellec-8 -rom
monitor=monitor.bin`, ma quei file restano esterni al repository finche'
provenance e licenze non saranno chiare.

I test di integrazione creano invece una ROM temporanea di tre byte:

```text
41 51 00
```

La ROM esegue `INP 0`, `OUT 8`, `HLT` e verifica il percorso completo tra file
locale, slot `test`, bus callback e trace I/O.

I file qui saranno aggiunti solo quando ci saranno vettori di test o programmi
di esempio stabili da versionare.
