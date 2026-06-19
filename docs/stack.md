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

Gli helper interni mascherano gli indirizzi a 14 bit e `SP` a 3 bit. `SP`
punta allo slot che contiene il PC corrente:

- durante il fetch, `setPC` aggiorna sia `PC` sia `Stack[SP]`
- `CALL` e `RST` avanzano `SP`, poi scrivono nello slot nuovo il target
- lo slot precedente conserva l'indirizzo di ritorno gia' avanzato dal fetch
- `RET` arretra `SP` e ripristina il PC dallo slot raggiunto

Questa scelta riflette la struttura interna dell'8008: 8 voci fisiche, ma 7
livelli utili di ritorno quando una voce e' occupata dal PC corrente.

---

## Implementato ora

- Campi stack e stack pointer nello stato CPU.
- Reset a zero dello stack.
- Helper primitivi per mascherare slot, SP e indirizzi.
- Mirroring del PC corrente in `Stack[SP]`.
- Push/pop impliciti per `CALL`, `RET` e `RST`.
- Overflow ciclico senza fault dopo l'ottavo livello fisico.
- Test sul mascheramento a 14 bit, profondita' 7 e overflow ciclico.

---

## Da implementare

- Integrazione con la futura logica di interrupt/jam instruction.
