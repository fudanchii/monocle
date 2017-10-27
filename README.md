Monocle
---

Another container based build tool.  
Useful to scope build / test only to related changes.

Example:
---

Given repository with folder `A` and `B`, and `build.yml` on each respective folders, and a commit with changes only at `B` folder. Monocle will detect changes and run build only at `B` folder


!!!
---
Monocle is a build tool, but it's also aimed to be on par with docker-compose to help in dev environment.  
It still lacks of multi-container / service containers support though. Also it's still has subpar UX for remote docker.

released under MIT.
