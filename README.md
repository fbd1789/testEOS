# eos-tester

**eos-tester** est un outil en Go permettant d'exécuter des tests de validation automatisés sur des équipements Arista EOS via eAPI, en s'appuyant sur :
- des définitions de tests en YAML (`config/test.yml`)
- des réponses JSON de commandes `show` EOS
- des templates Go (`text/template`) pour évaluer les résultats

---

## 📦 Structure du projet

```
eos-tester/
├── api/                   # Client eAPI Arista (authentification, exécution de commande)
│   └── client.go
├── config/
│   └── test.yml           # Fichier de définition des tests à exécuter
├── engine/                # Moteur de test (exécution, validation, logique de rendu)
│   ├── runner.go
│   ├── validator.go
│   └── version.go
├── parser/
│   └── loader.go          # Chargement des tests YAML
├── templates/             # Templates Go pour chaque type de test (comparaison JSON)
│   ├── version_check.tmpl
│   ├── power_state_check.tmpl
│   ├── tempsensor_check.tmpl
│   ├── fan_status_check.tmpl
│   ├── fan_speed_check.tmpl
│   ├── timezone_check.tmpl
│   ├── ntp_check.tmpl
│   └── ntp_status_check.tmpl
├── hosts.txt              # Liste des équipements à tester (1 IP/DNS par ligne)
├── go.mod / go.sum        # Dépendances du projet
└── main.go                # Point d'entrée CLI
```

---

## ▶️ Lancer des tests

### ✔️ Tous les tests sur tous les hôtes :
```bash
go run main.go -hosts hosts.txt -user <login> -pass <motdepasse>
```

### ✔️ Un seul test ciblé :
```bash
go run main.go -hosts hosts.txt -user <login> -pass <motdepasse> -s "Check Timezone"
```

---

## 🧪 Exemple de test (`test.yml`)

```yaml
tests:
  - name: Check Timezone
    command: "show clock"
    template: "timezone_check.tmpl"
    vars:
      timezone: "Europe/Paris"

  - name: Check NTP Servers Present
    command: "show clock"
    template: "ntp_check.tmpl"
    vars:
      ntp_servers:
        - "129.250.35.250"
        - "ntp1.example.org"
```

---

## 🧰 Templates

Les templates utilisent la syntaxe Go `text/template` avec [Sprig](https://masterminds.github.io/sprig/) intégré. Exemple de logique :

```gotmpl
{{- $state := index .result "status" }}
{{- if eq $state "synchronised" }}
OK
{{- else }}
FAIL - status incorrect: {{ $state }}
{{- end }}
```

---

## 🧩 Ajouter un nouveau test

1. Écrire un template `.tmpl` dans `templates/`
2. Ajouter une entrée dans `config/test.yml`
3. Relancer l’outil

---

## 🔐 Connexion

La connexion aux équipements EOS se fait via HTTPS (eAPI). Le fichier `api/client.go` gère :
- Authentification basique
- Transport TLS sans vérification (test local)
