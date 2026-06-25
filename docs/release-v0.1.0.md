# Checklist release v0.1.0

La prima release pubblica dell'emulatore fotografa il core instruction-accurate
e le milestone 0-17 e 20-25. Le mappe storiche e le ROM delle milestone 18-19
restano esplicitamente fuori dallo scope.

## Gate obbligatori

- [ ] Worktree pulito e modifiche di release revisionate.
- [ ] `gofmt -l .` non produce output.
- [ ] `go vet ./...` termina senza errori.
- [ ] `go test -count=1 ./...` termina senza errori.
- [ ] `go run ./cmd/retronet-8008 -conformance` termina con tutti i casi verdi.
- [ ] Il workflow di integrazione dell'ecosistema assembla ed esegue la demo
      8008 tramite `retronet-asm`.
- [ ] README, roadmap e limiti noti descrivono lo stesso perimetro.
- [ ] Il commit candidato usa dipendenze RetroNet pubblicate e riproducibili.

## Limiti dichiarati della release

- Nessuna ROM storica distribuita.
- Profili SCELBI/Intellec conservativi, senza mappe memoria/I/O inventate.
- Timing per istruzione e ciclo macchina, non transizioni elettriche dei pin.
- Conformance esaustiva interna, ma non ancora differenziale contro un secondo
  emulatore indipendente.

## Pubblicazione

Dopo il completamento dei gate, creare il tag annotato sul commit candidato:

```bash
git tag -a v0.1.0 -m "retronet-8008 v0.1.0"
git push origin main
git push origin v0.1.0
```

Il tag non va creato su un worktree sporco o prima che la CI del commit
candidato sia verde.
