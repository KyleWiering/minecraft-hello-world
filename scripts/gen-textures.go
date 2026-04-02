// gen-textures.go — generates cat-themed replacement textures for primary Minecraft mobs
// and the helloworld block texture.
//
// Run from the project root with:
//
//	go run ./scripts/gen-textures.go
//
// Outputs PNG files into resource_pack/textures/ at the same paths that
// Minecraft Bedrock expects.
//
// Mobs and their cat themes
// ─────────────────────────
//
//	zombie        → Orange tabby        (64 × 64)
//	skeleton      → Cream/white cat     (64 × 64)
//	creeper       → Black cat           (64 × 32)
//	spider        → Siamese             (64 × 32)
//	enderman      → Persian gray        (64 × 32)
//	cow           → Tuxedo              (64 × 32)
//	pig           → Calico              (64 × 32)
//	chicken       → White cat           (64 × 32)
//	sheep body    → Ragdoll tan         (64 × 32)
//	sheep wool    → Ragdoll cream tabby (64 × 32)
//
//	helloworld block → Green with "HW" pixel art  (16 × 16)
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// savePNG writes img as a PNG to relPath, which is relative to resource_pack/textures/.
// The script must be run from the project root.
func savePNG(relPath string, img image.Image) {
	full := filepath.Join("resource_pack", "textures", relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		panic(err)
	}
	f, err := os.Create(full)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
	b := img.Bounds()
	fmt.Printf("  ✔  textures/%s  (%d×%d)\n", relPath, b.Max.X-b.Min.X, b.Max.Y-b.Min.Y)
}

// newRGBA returns a new RGBA image of the given dimensions.
func newRGBA(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

// fillSolid paints every pixel of img with c.
func fillSolid(img *image.RGBA, c color.RGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.SetRGBA(x, y, c)
		}
	}
}

// fillTabby paints horizontal tabby stripes on img, alternating between base and stripe.
func fillTabby(img *image.RGBA, base, stripe color.RGBA, stripeWidth, offset int) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if ((y+offset)/stripeWidth)%2 == 1 {
				img.SetRGBA(x, y, stripe)
			} else {
				img.SetRGBA(x, y, base)
			}
		}
	}
}

// fillSiamese paints a Siamese colour-point pattern: light body, dark points fading in
// from the edges.
func fillSiamese(img *image.RGBA, body, points color.RGBA, blendDepth int) {
	b := img.Bounds()
	w := b.Max.X
	h := b.Max.Y
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			edgeDist := x
			if y < edgeDist {
				edgeDist = y
			}
			if w-1-x < edgeDist {
				edgeDist = w - 1 - x
			}
			if h-1-y < edgeDist {
				edgeDist = h - 1 - y
			}
			var c color.RGBA
			if edgeDist < blendDepth {
				t := 1.0 - float64(edgeDist)/float64(blendDepth)
				c = color.RGBA{
					R: uint8(float64(body.R) + t*float64(int(points.R)-int(body.R))),
					G: uint8(float64(body.G) + t*float64(int(points.G)-int(body.G))),
					B: uint8(float64(body.B) + t*float64(int(points.B)-int(body.B))),
					A: 255,
				}
			} else {
				c = body
			}
			img.SetRGBA(x, y, c)
		}
	}
}

// fillTuxedo paints a black-and-white tuxedo pattern: black upper-left, white lower-right.
func fillTuxedo(img *image.RGBA) {
	b := img.Bounds()
	w := float64(b.Max.X)
	h := float64(b.Max.Y)
	black := color.RGBA{24, 24, 24, 255}
	white := color.RGBA{240, 240, 240, 255}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			fx := float64(x) / w
			fy := float64(y) / h
			if (fx < 0.35 && fy < 0.55) || (fx > 0.65 && fy > 0.45) {
				img.SetRGBA(x, y, black)
			} else {
				img.SetRGBA(x, y, white)
			}
		}
	}
}

// fillCalico paints orange/white/black calico rectangular patches.
func fillCalico(img *image.RGBA) {
	type patch struct{ x0, y0, x1, y1 float64; c color.RGBA }
	patches := []patch{
		{0.00, 0.00, 0.40, 0.50, color.RGBA{228, 96, 18, 255}},  // orange
		{0.40, 0.00, 1.00, 0.40, color.RGBA{240, 240, 240, 255}}, // white
		{0.00, 0.50, 0.50, 1.00, color.RGBA{240, 240, 240, 255}}, // white
		{0.50, 0.45, 1.00, 0.70, color.RGBA{28, 28, 28, 255}},   // black
		{0.50, 0.70, 1.00, 1.00, color.RGBA{228, 96, 18, 255}},  // orange
	}
	b := img.Bounds()
	w := float64(b.Max.X)
	h := float64(b.Max.Y)
	defColor := color.RGBA{240, 240, 240, 255}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			fx := float64(x) / w
			fy := float64(y) / h
			c := defColor
			for _, p := range patches {
				if fx >= p.x0 && fx < p.x1 && fy >= p.y0 && fy < p.y1 {
					c = p.c
					break
				}
			}
			img.SetRGBA(x, y, c)
		}
	}
}

// ── Mob texture generators ───────────────────────────────────────────────────

// genZombie generates the orange tabby zombie texture (64×64).
func genZombie() {
	img := newRGBA(64, 64)
	fillTabby(img, color.RGBA{224, 120, 48, 255}, color.RGBA{160, 72, 16, 255}, 4, 0)
	savePNG("entity/zombie/zombie.png", img)
}

// genSkeleton generates the cream/white skeleton texture (64×64).
func genSkeleton() {
	img := newRGBA(64, 64)
	fillTabby(img, color.RGBA{240, 232, 208, 255}, color.RGBA{200, 185, 160, 255}, 6, 0)
	savePNG("entity/skeleton/skeleton.png", img)
}

// genCreeper generates the black cat creeper texture (64×32).
func genCreeper() {
	img := newRGBA(64, 32)
	fillSolid(img, color.RGBA{22, 22, 22, 255})
	gold := color.RGBA{255, 210, 0, 255}
	// Two golden eye patches on the face UV region.
	for y := 4; y <= 8; y++ {
		for x := 10; x <= 17; x++ {
			img.SetRGBA(x, y, gold)
		}
		for x := 42; x <= 49; x++ {
			img.SetRGBA(x, y, gold)
		}
	}
	savePNG("entity/creeper/creeper.png", img)
}

// genSpider generates the Siamese spider texture (64×32).
func genSpider() {
	img := newRGBA(64, 32)
	fillSiamese(img, color.RGBA{245, 222, 179, 255}, color.RGBA{61, 31, 0, 255}, 10)
	savePNG("entity/spider/spider.png", img)
}

// genEnderman generates the Persian gray enderman texture (64×32).
func genEnderman() {
	img := newRGBA(64, 32)
	fillTabby(img, color.RGBA{88, 78, 108, 255}, color.RGBA{56, 48, 72, 255}, 4, 0)
	savePNG("entity/enderman/enderman.png", img)
}

// genCow generates the tuxedo cow texture (64×32).
func genCow() {
	img := newRGBA(64, 32)
	fillTuxedo(img)
	savePNG("entity/cow/cow.png", img)
}

// genPig generates the calico pig texture (64×32).
func genPig() {
	img := newRGBA(64, 32)
	fillCalico(img)
	savePNG("entity/pig/pig.png", img)
}

// genChicken generates the white cat chicken texture (64×32).
func genChicken() {
	img := newRGBA(64, 32)
	fillTabby(img, color.RGBA{250, 248, 240, 255}, color.RGBA{210, 200, 185, 255}, 6, 0)
	savePNG("entity/chicken.png", img)
}

// genSheep generates the Ragdoll sheep body and wool textures (64×32 each).
func genSheep() {
	body := newRGBA(64, 32)
	fillSolid(body, color.RGBA{160, 148, 132, 255})
	savePNG("entity/sheep/sheep.png", body)

	wool := newRGBA(64, 32)
	fillTabby(wool, color.RGBA{236, 216, 180, 255}, color.RGBA{200, 175, 135, 255}, 5, 0)
	savePNG("entity/sheep/sheep_fur.png", wool)
}

// ── Block texture generator ──────────────────────────────────────────────────

// genHelloWorld generates the helloworld block texture (16×16): a green block
// with pixel-art "HW" letters in yellow.
func genHelloWorld() {
	img := newRGBA(16, 16)

	bg := color.RGBA{45, 120, 45, 255}    // medium green
	border := color.RGBA{20, 80, 20, 255} // dark green border
	letter := color.RGBA{255, 220, 0, 255} // yellow letters

	fillSolid(img, bg)

	// 1-pixel dark border
	for i := 0; i < 16; i++ {
		img.SetRGBA(i, 0, border)
		img.SetRGBA(i, 15, border)
		img.SetRGBA(0, i, border)
		img.SetRGBA(15, i, border)
	}

	// "H" at (2, 5), 5 wide × 5 tall
	//  X . . . X
	//  X . . . X
	//  X X X X X
	//  X . . . X
	//  X . . . X
	hx, hy := 2, 5
	for row := 0; row < 5; row++ {
		img.SetRGBA(hx+0, hy+row, letter)
		img.SetRGBA(hx+4, hy+row, letter)
	}
	for col := 1; col < 4; col++ {
		img.SetRGBA(hx+col, hy+2, letter)
	}

	// "W" at (9, 5), 5 wide × 5 tall
	//  X . . . X
	//  X . . . X
	//  X . X . X
	//  . X . X .
	//  . . X . .
	wx, wy := 9, 5
	img.SetRGBA(wx+0, wy+0, letter)
	img.SetRGBA(wx+4, wy+0, letter)
	img.SetRGBA(wx+0, wy+1, letter)
	img.SetRGBA(wx+4, wy+1, letter)
	img.SetRGBA(wx+0, wy+2, letter)
	img.SetRGBA(wx+2, wy+2, letter)
	img.SetRGBA(wx+4, wy+2, letter)
	img.SetRGBA(wx+1, wy+3, letter)
	img.SetRGBA(wx+3, wy+3, letter)
	img.SetRGBA(wx+2, wy+4, letter)

	savePNG("blocks/helloworld.png", img)
}

// ── Entry point ──────────────────────────────────────────────────────────────

func main() {
	fmt.Println("Generating CatMob Madness textures…\n")
	genZombie()
	genSkeleton()
	genCreeper()
	genSpider()
	genEnderman()
	genCow()
	genPig()
	genChicken()
	genSheep()

	fmt.Println("\nGenerating helloworld block texture…\n")
	genHelloWorld()

	fmt.Println("\nDone ✅  — edit the PNGs in resource_pack/textures/ to customise each skin.")
}
