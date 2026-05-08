# Servizio **`iot-nonna-core`**

## Scopo generale

**`iot-nonna-core`** è il servizio centrale di dominio del sistema IoT “nonna”.  
Ha la responsabilità di:

- **gestire il modello business IoT** (stanze, dispositivi, sensori, tipologie)
- **esporre i dati raccolti** in forma chiara e consumabile
- **nascondere la complessità del database** al resto del sistema
- **fungere da punto di verità unico** per configurazione e lettura dei dati

Non è un servizio di ingestione dati né di controllo attivo, ma **il cuore informativo e configurativo** su cui tutti gli altri servizi si appoggiano.

---

## Responsabilità del servizio

### ✅ Cosa **DEVE** fare
1. **Provisioning e configurazione**
    - Creazione e gestione di:
        - rooms
        - device_type
        - sensor_type
        - devices
        - associazioni device ↔ sensori
    - Applicazione di:
        - validazioni
        - vincoli logici
        - regole di coerenza (es. un device deve avere un type valido)
2. **Accesso semplificato ai dati**
    - Esporre il contenuto del database in forma:
        - leggibile
        - strutturata
        - stabile nel tempo
    - Eliminare la necessità di usare pgAdmin o query manuali
3. **Lettura dei dati storici**
    - Fornire endpoint pensati per:
        - dashboard
        - grafici
        - analisi temporali
    - Incapsulare:
        - filtri temporali
        - aggregazioni
        - downsampling
    - Risultati pronti per il frontend (non “raw SQL”)
4. **Stabilità del dominio**
    - Essere il **contratto stabile** tra DB e consumer
    - Permettere evoluzione dello schema DB senza rompere il frontend

---

## Scelte tecnologiche
- **Linguaggio**: Go
- **Stile API**: REST
- **Router HTTP**: `chi`
- **Database**: PostgreSQL
- **Accesso DB**: driver nativo (`pgx`)
- **Separazione ruoli DB**:
    - `iot-nonna-ingest`: lettura metadata + scrittura readings
    - `iot-nonna-core`: gestione tabelle business + lettura

### Perché Go + chi
- Go garantisce:
    - performance prevedibili
    - semplicità strutturale
    - eccellente supporto a servizi backend di dominio
- `chi`:
    - minimalista
    - basato su `net/http`
    - non introduce magia o pattern artificiali
    - favorisce architetture pulite e mantenibili nel tempo

---

## Struttura logica del servizio

Il servizio è **stratificato**, non “a handler grassi”.

Concetti chiave (indipendenti dal codice):

1. **HTTP layer**
    
    - parsing request
    - status code
    - serializzazione JSON
2. **Domain / Service layer**
    
    - regole di business
    - validazioni logiche
    - orchestrazione delle operazioni
3. **Persistence layer**
    
    - query SQL
    - accesso al DB
    - nessuna logica di dominio

Ogni layer ha responsabilità **chiare e non sovrapposte**.

---

## Migration e Seeding

### Migration (obbligatorie)

Le **migration** servono a:

- versionare lo schema del database
- rendere l’evoluzione del DB riproducibile e controllata

Ogni modifica strutturale (tabelle, colonne, indici) è:

- descritta in un file SQL numerato
- applicata in ordine

Questo elimina:

- modifiche manuali
- differenze tra ambienti
- “non so come è nato questo campo”

Le migration fanno **parte del servizio**, non di un servizio separato.

---

### Seeding (dati iniziali)

Il **seeding** serve a:

- popolare il DB con dati iniziali
- rendere immediato avvio e test
- evitare configurazioni manuali

Esempi:

- tipi di sensori noti
- tipi di dispositivi noti
- configurazione di test

Il seeding può essere:

- eseguito manualmente
- o all’avvio in ambienti non produttivi