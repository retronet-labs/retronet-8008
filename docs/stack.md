# Stack Interno

Lo stack dell'Intel 8008 non e' uno stack dati in memoria: e' una struttura
interna al processore usata per program counter e indirizzi di ritorno.

---

## Cosa rappresenta

Lo stack conserva indirizzi a 14 bit. Ha 8 voci fisiche, ma una rappresenta il
PC corrente; per questo i livelli utili di annidamento per `CALL` e `RST` sono 7.

---

## Come funziona nell'8008

Il puntatore stack e' interno, a 3 bit, e non e' visibile al programma. Le
istruzioni `CALL` e `RST` salvano l'indirizzo di ritorno; `RET` lo ripristina.
Se si supera la capacita', il puntatore ricicla e una voce precedente viene
sovrascritta senza fault.

---

## Come e' modellato nel progetto

`CPU8008` contiene:

- `Stack [8]uint16`
- `SP uint8`

Gli helper interni mascherano gli indirizzi a 14 bit e `SP` a 3 bit. La semantica
esatta di push/pop per `CALL`, `RET` e `RST` sara' fissata nella milestone sul
controllo di flusso.

---

## Implementato ora

- Campi stack e stack pointer nello stato CPU.
- Reset a zero dello stack.
- Helper primitivi per mascherare slot, SP e indirizzi.
- Test sul mascheramento a 14 bit.

---

## Da implementare

- Push/pop effettivi.
- `CALL`, `RET` e `RST`.
- Test di profondita' 1, profondita' 7 e overflow ciclico.
