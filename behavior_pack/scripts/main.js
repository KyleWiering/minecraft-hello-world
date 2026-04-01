import { world } from "@minecraft/server";

const GREETING = "§a§lHello World!§r §7Welcome to Minecraft — enjoy your adventure!";

// Greet each player the first time they spawn in the world.
world.afterEvents.playerSpawn.subscribe((event) => {
  if (event.initialSpawn) {
    event.player.sendMessage(GREETING);
  }
});
