import { world } from "@minecraft/server";

const GREETING =
  "§6§lCatMob Madness v1§r §eis active! §7Every mob has been turned into a cat. Meow!";

// Greet each player the first time they spawn in the world.
world.afterEvents.playerSpawn.subscribe((event) => {
  if (event.initialSpawn) {
    event.player.sendMessage(GREETING);
  }
});
