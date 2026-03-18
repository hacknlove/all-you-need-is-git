---
description: Releases a new version
agent: build
---

Desplegar la versión $1

Si no se ha indicado versión, analizar la version actual en @docs/public/releases/latest.txt y analizar los mensajes de commits desde ese tag para decidir si se trata de un mayor, un minor o un patch

El número de versión ha de seguir el formato vX.y.Z, tanto para el tag como para el pointer docs/public/releases/latest.txt

Una vez el número de versión está claro:

Ejecutar: `./scripts/release.sh vX.Y.Z`
