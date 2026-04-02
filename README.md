# CatMob Madness v1

A **CatMob Madness** add-on for **Minecraft Bedrock Edition** that:

- Replaces the textures of every primary overworld mob with a cat-themed skin.
- Sends every player a CatMob Madness announcement the first time they spawn.
- Works locally on Windows 10/11 for quick testing.
- Loads in **Realms** and any Bedrock world (including Xbox).
- Can be packaged and submitted to the **Xbox / Minecraft Marketplace**.

## Cat themes at a glance

| Mob       | Cat variety         | Texture size |
|-----------|---------------------|--------------|
| Zombie    | Orange tabby        | 64 × 64      |
| Skeleton  | Cream/white cat     | 64 × 64      |
| Creeper   | Black cat           | 64 × 32      |
| Spider    | Siamese             | 64 × 32      |
| Enderman  | Persian gray        | 64 × 32      |
| Cow       | Tuxedo              | 64 × 32      |
| Pig       | Calico              | 64 × 32      |
| Chicken   | White cat           | 64 × 32      |
| Sheep     | Ragdoll (body + wool) | 64 × 32    |

---

## Customising the skins

Each mob's texture lives in `resource_pack/textures/entity/` at the same path that
Minecraft uses for its built-in textures.  Open any PNG in your favourite pixel-art
editor (e.g. **Aseprite**, **GIMP**, or **Paint.NET**) and repaint it — no code
changes are required.

To regenerate the stock cat textures from scratch (e.g. after accidentally deleting
one), run:

```bash
python3 scripts/gen-textures.py
```

> **Tip:** The generator script (`scripts/gen-textures.py`) is heavily commented and
> uses only Python's built-in `struct` / `zlib` modules — no third-party packages
> needed.

---

## Repository structure

```
minecraft-hello-world/
├── behavior_pack/              # Game-logic pack
│   ├── manifest.json
│   └── scripts/
│       └── main.js             # CatMob Madness spawn-announcement script
├── resource_pack/              # Asset pack
│   ├── manifest.json
│   └── textures/
│       └── entity/             # Cat-themed mob textures (PNG)
│           ├── zombie/zombie.png
│           ├── skeleton/skeleton.png
│           ├── creeper/creeper.png
│           ├── spider/spider.png
│           ├── enderman/enderman.png
│           ├── cow/cow.png
│           ├── pig/pig.png
│           ├── chicken.png
│           └── sheep/
│               ├── sheep.png
│               └── sheep_fur.png
├── scripts/
│   ├── pack.js                 # Build script → produces dist/catmob-madness.mcaddon
│   └── gen-textures.py         # Regenerates stock cat textures
├── package.json
└── .gitignore
```

---

## Local testing (Windows)

### 1 — Install the packs

Copy both folders to the Minecraft `com.mojang` directory:

```powershell
$comMojang = "$env:LOCALAPPDATA\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang"

Copy-Item -Recurse -Force .\behavior_pack "$comMojang\behavior_packs\catmob-madness-bp"
Copy-Item -Recurse -Force .\resource_pack "$comMojang\resource_packs\catmob-madness-rp"
```

### 2 — Enable the packs on a world

1. Open Minecraft and create (or edit) a world.
2. Go to **Add-Ons** → **Behavior Packs**, find *CatMob Madness Behavior Pack*, and activate it.
3. Go to **Add-Ons** → **Resource Packs**, find *CatMob Madness Resource Pack*, and activate it.
4. Launch the world — you will see the CatMob Madness announcement the moment you spawn, and every mob will wear a cat skin.

---

## Loading in Realms

1. Open the Realm settings in Minecraft.
2. Navigate to **Add-Ons** and enable both the behavior pack and resource pack as shown above.
3. Apply the changes and re-join the Realm — the announcement appears for every new player and all mobs display cat textures.

---

## Packaging into a distributable `.mcaddon`

A `.mcaddon` is a single zip file that installs both packs at once.

```bash
npm install       # install dev dependency (archiver)
npm run pack      # produces dist/catmob-madness.mcaddon
```

Double-clicking `catmob-madness.mcaddon` on Windows automatically imports both packs into Minecraft.

---

## Submitting to the Xbox One / Minecraft Marketplace

> The Marketplace lets players download your add-on from inside Minecraft on Xbox, Windows, iOS, Android, and more.

### Prerequisites

- A **Microsoft Partner Center** account — sign up at <https://partner.microsoft.com>.
- Enroll in the **Minecraft Marketplace Partner Program** — apply at <https://aka.ms/MCMarketplacePartner>.
- Review the [Minecraft Creator Portal](https://learn.microsoft.com/en-us/minecraft/creator/) for content guidelines and technical requirements.

### Submission steps

1. Ensure your `manifest.json` files use unique UUIDs (update the placeholder UUIDs in this repo before publishing).
2. Run `npm run pack` to create `dist/catmob-madness.mcaddon`.
3. Log in to [Partner Center](https://partner.microsoft.com/dashboard) and create a new **Minecraft Add-On** product.
4. Upload `catmob-madness.mcaddon` as the content package.
5. Fill in the store listing details (title, description, screenshots, pricing).
6. Submit for Microsoft certification — the review process typically takes a few business days.
7. Once approved, the add-on goes live on the Xbox One Store and all Minecraft Marketplace storefronts.

---

## Minecraft Script API version compatibility

| `@minecraft/server` version | Minecraft Bedrock version |
|-----------------------------|---------------------------|
| 1.5.0                       | 1.20.x and later          |

The `min_engine_version` in both manifests is set to `[1, 20, 0]`.
