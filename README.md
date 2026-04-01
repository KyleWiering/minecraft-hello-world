# minecraft-hello-world

A **Hello World** add-on for **Minecraft Bedrock Edition** that:

- Sends every player a greeting message the first time they spawn.
- Works locally on Windows 10/11 for quick testing.
- Loads in **Realms** and any Bedrock world (including Xbox).
- Can be packaged and submitted to the **Xbox / Minecraft Marketplace**.

---

## Repository structure

```
minecraft-hello-world/
├── behavior_pack/          # Game-logic pack
│   ├── manifest.json
│   └── scripts/
│       └── main.js         # Hello-world player-spawn script
├── resource_pack/          # Asset pack (companion to behavior pack)
│   └── manifest.json
├── scripts/
│   └── pack.js             # Build script → produces dist/hello-world.mcaddon
├── package.json
└── .gitignore
```

---

## Local testing (Windows)

### 1 — Install the packs

Copy both folders to the Minecraft `com.mojang` directory:

```powershell
$comMojang = "$env:LOCALAPPDATA\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang"

Copy-Item -Recurse -Force .\behavior_pack "$comMojang\behavior_packs\hello-world-bp"
Copy-Item -Recurse -Force .\resource_pack "$comMojang\resource_packs\hello-world-rp"
```

### 2 — Enable the packs on a world

1. Open Minecraft and create (or edit) a world.
2. Go to **Add-Ons** → **Behavior Packs**, find *Hello World Behavior Pack*, and activate it.
3. Go to **Add-Ons** → **Resource Packs**, find *Hello World Resource Pack*, and activate it.
4. Launch the world — you will see the green greeting message the moment you spawn.

---

## Loading in Realms

1. Open the Realm settings in Minecraft.
2. Navigate to **Add-Ons** and enable both the behavior pack and resource pack as shown above.
3. Apply the changes and re-join the Realm — the greeting will appear for every new player.

---

## Packaging into a distributable `.mcaddon`

A `.mcaddon` is a single zip file that installs both packs at once.

```bash
npm install       # install dev dependency (archiver)
npm run pack      # produces dist/hello-world.mcaddon
```

Double-clicking `hello-world.mcaddon` on Windows automatically imports both packs into Minecraft.

---

## Submitting to the Xbox One / Minecraft Marketplace

> The Marketplace lets players download your add-on from inside Minecraft on Xbox, Windows, iOS, Android, and more.

### Prerequisites

- A **Microsoft Partner Center** account — sign up at <https://partner.microsoft.com>.
- Enroll in the **Minecraft Marketplace Partner Program** — apply at <https://aka.ms/MCMarketplacePartner>.
- Review the [Minecraft Creator Portal](https://learn.microsoft.com/en-us/minecraft/creator/) for content guidelines and technical requirements.

### Submission steps

1. Ensure your `manifest.json` files use unique UUIDs (update the placeholder UUIDs in this repo before publishing).
2. Run `npm run pack` to create `dist/hello-world.mcaddon`.
3. Log in to [Partner Center](https://partner.microsoft.com/dashboard) and create a new **Minecraft Add-On** product.
4. Upload `hello-world.mcaddon` as the content package.
5. Fill in the store listing details (title, description, screenshots, pricing).
6. Submit for Microsoft certification — the review process typically takes a few business days.
7. Once approved, the add-on goes live on the Xbox One Store and all Minecraft Marketplace storefronts.

---

## Minecraft Script API version compatibility

| `@minecraft/server` version | Minecraft Bedrock version |
|-----------------------------|---------------------------|
| 1.5.0                       | 1.20.x and later          |

The `min_engine_version` in both manifests is set to `[1, 20, 0]`.
