package gameDoc

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"
)

var walls []pixel.Rect
var imd = imdraw.New(nil)
var imd1 = imdraw.New(nil)
var imd2 = imdraw.New(nil)
var imd3 = imdraw.New(nil)
var imd4 = imdraw.New(nil)
var imd5 = imdraw.New(nil)
var imd6 = imdraw.New(nil)
var imd7 = imdraw.New(nil)
var imd8 = imdraw.New(nil)
var badGuys []BadGuy
var xLoc = []float64{50, 100, 160, 220, 280, 340, 400, 460, 540, 600, 660, 720}
var yLoc = []float64{120, 170, 240, 310, 400, 480, 540, 560}
var lives = 0
var score = 0
var level = 1

//GameLevel struct will define the level properties
type GameLevel struct{
	availableX int
	availableY int
}

//Stuct Player will define the properties for the Player
//Contains the sprite for the player
// Contains the vector for the positioning of the player
type Player struct{
	Image     *pixel.Sprite
	Direction string
	Body      pixel.Vec
	Shot      Bullet
	ShotFired bool
	Lives     int
	Score     int
}


//Struct BadGuy will define the properties for the bad guys
//Contains the Sprite of the BadGuy
//Contains the Vector fot the positioning of the BadGuy
type BadGuy struct{
	Image *pixel.Sprite
	Body pixel.Vec
	Shot Bullet
	ShotFired bool
	Destination pixel.Vec
	Direction string
	Alive bool
	Count int
}

//Struct Bullet will define the properties for the badGuy and the player Bullet
// Contains the sprite for the Bullet
// Contains the direction and the vector of the positioning of the Bullet
type Bullet struct{
	Image *pixel.Sprite
	Body pixel.Vec
	Hit bool
	Direction string
}

//The main function calls the Run method of the pixel package and passes Start function as parameters.
func Main() {
	pixelgl.Run(Start)
}

//Code resued from the Github Tutorial for pixel and from class demo
//The Start function will define the window properties and will setup the call for the first level of the game.
//Start function call the Run function setting up the properies and the Level Properties
func Start() {
	cfg := pixelgl.WindowConfig{
		Title:  "Berzerk",
		Bounds: pixel.R(0, 0, 800, 640),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Dimgray)
	SetLevelOneWall()
	Run(win, NewPlayer(), GameLevel{11,7})
}

//The function Run takes an object of pixelgl Window, a Player Struct and a GameLevel Struct.
//The Run function will characteristics of the level and create bad guys as per the GameLevel.
//The Run function is responsible for drawing the elements in the game.
//The function handles the operation within the window and updates the window during changes.
func Run(win *pixelgl.Window, player Player, gameLevel GameLevel) {
	badGuyBullet := NewBullet()
	badGuyShot := false
	win.Clear(colornames.Dimgray)
	lives = player.Lives

	for i:= 0; i< (4 +(2* level)); i++{
		badGuy := NewBadGuy()
		badGuy.Body.X = xLoc[rand.Intn(gameLevel.availableX)]
		badGuy.Body.Y = yLoc[rand.Intn(gameLevel.availableY)]
		badGuy.Count = rand.Intn(100)
		badGuys = append(badGuys, badGuy)
	}
	for !win.Closed(){
		win.Clear(colornames.Dimgray)
		imd.Draw(win)
		imd1.Draw(win)
		imd2.Draw(win)
		imd3.Draw(win)
		imd4.Draw(win)
		imd5.Draw(win)
		imd6.Draw(win)
		imd7.Draw(win)
		imd8.Draw(win)

		DrawScore(win, player.Score)
		DrawLevel(win)

		for i, badguy := range badGuys{
			if badguy.Alive == true{
				badguy.Image.Draw(win,pixel.IM.Moved(pixel.V(badguy.Body.X, badguy.Body.Y)))
				badguy, player = UpdateLevelBadGuys(badguy, player, i)
				if badguy.Count % 200 == 0 {
					if badGuyShot == false{
						badGuyShot = true
						PlaySoundEffects("assets/badGuyBullet.mp3")
						badGuyBullet.Body.X = badguy.Body.X
						badGuyBullet.Body.Y = badguy.Body.Y
						badGuyBullet.Direction = badguy.Direction
						badGuyBulletPic, _ := LoadPicture("assets/bullet.png")
						badGuyBullet.Image.Set(badGuyBulletPic, badGuyBulletPic.Bounds())
						badGuyBullet.Image.Draw(win, pixel.IM.Moved(pixel.V(badGuyBullet.Body.X, badGuyBullet.Body.Y)))
					}
				}
			}
			if BulletHit(player.Shot, badguy){
				player.ShotFired = false
				player.Shot.Body = pixel.Vec{}
				player.Shot.Image.Set(nil, pixel.Rect{})
			}
			badGuys = append(badGuys[:i], badguy)
		}

		if badGuyShot == true{
			if badGuyBullet.Direction == "Left" {
				badGuyBullet.Body.X -= 4
			}
			if badGuyBullet.Direction == "Right" {
				badGuyBullet.Body.X += 4
			}
			if badGuyBullet.Direction == "Up"  {
				badGuyBullet.Body.Y += 4
			}
			if badGuyBullet.Direction == "Down"  {
				badGuyBullet.Body.Y -= 4
			}
			badGuyBullet.Image.Draw(win, pixel.IM.Moved(pixel.V(badGuyBullet.Body.X, badGuyBullet.Body.Y)))
		}
		if BadGuysBulletHitWall(badGuyBullet){
			badGuyShot = false
			badGuyBullet.Body = pixel.Vec{}
			badGuyBullet.Image.Set(nil, pixel.Rect{})
		}
		if BadGuyHitPlayer(badGuyBullet, player){
			badGuyShot = false
			PlayerDead(win, player)
			score := player.Score
			player = NewPlayer()
			lives --
			PlaySoundEffects("assets/playerLoseLife.mp3")
			player.Lives = lives
			player.Score = score
			badGuyBullet.Body = pixel.Vec{}
			badGuyBullet.Image.Set(nil, pixel.Rect{})
		}
		if BadGuysAllDead(){
			if level == 4{
				GameWon(win, player)
			}
			if player.Body.Y < 40.0 || player.Body.Y > 635.0 {
				level++
				PlaySoundEffects("assets/levelUp.mp3")
				ChangeLevels(win, player)
			}
		}
		player = UpdateLevelPlayer(win, player)
		player.Image.Draw(win, pixel.IM.Moved(pixel.V(player.Body.X, player.Body.Y)))
		DrawLives(win, player)
		win.Update()
	}
}

// The function will return a struct of Player
// The function will set a sprite for the player and also defines the vector that will set the position of the player.
func NewPlayer() Player{
	pic , _ := LoadPicture("assets/player.png")
	return Player{
		pixel.NewSprite(pic, pic.Bounds()),
		"Right",
		pixel.Vec{40.0, 384.0},
		NewBullet(),
		false,
		3,
		0,
	}
}

//The function will return a struct of Bad Guy
// The function will set a sprite for the Bad Guy and defines the vector that will set the position of the Bad Guy
func NewBadGuy() BadGuy {
	pic , _ := LoadPicture("assets/bad_guy.png")
	return BadGuy{
		pixel.NewSprite(pic, pic.Bounds()),
		pixel.Vec{},
		Bullet{},
		false,
		pixel.Vec{},
		"",
		true,
		0,
	}
}

// The function will return a struct of Bad Guy
// The function will set a sprite for the Bullet and defines the vector that will set the position of the Bullet
func NewBullet() Bullet{
	pic, _ := LoadPicture("assets/bullet.png")
	return Bullet{
		pixel.NewSprite(pic, pic.Bounds()),
		pixel.Vec{},
		false,
		"",
	}
}

// Code from pixel Github page.
// The function LoadPicture takes a string path for the picture as argument and returns a picture
func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	img.Bounds();
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

// The function takes window and an integer (number of lives) as argument.
// The function will loop the number and draw score at the bottom Right corner.
func DrawScore(win *pixelgl.Window, score int){
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(700, 30), atlas)
	basicTxt.Color = color.Black
	fmt.Fprintln(basicTxt, "SCORE:", score)
	basicTxt.Draw(win, pixel.IM)
}

//The function takes pixel window as arguments and will draw the level
//The Level number will be drawn in the center bottom.
func DrawLevel(win *pixelgl.Window){
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(375, 30), atlas)
	basicTxt.Color = color.Black
	fmt.Fprintln(basicTxt, "LEVEL ", level)
	basicTxt.Draw(win, pixel.IM)
}

// The function takes pixel window and a struct of Player as arguments and returns a Struct of Player as output
// The function updates the field of the Player of the Game
// The function determines the direction of the player based on the key pressed.
// The function sets the bullet of the Player.
// The UpdateLevelPlayer function updates the score and lives of the player
func UpdateLevelPlayer(win *pixelgl.Window, player Player) Player{
	if IsPlayerDead(player){
		PlayerDead(win, player)
		score = player.Score
		player = NewPlayer()
		lives --
		PlaySoundEffects("assets/playerLoseLife.mp3")
		player.Lives = lives
		player.Score = score
	}
	if win.Pressed(pixelgl.KeyLeft) {
		if win.Pressed(pixelgl.KeyDown){
			player.Direction = "LeftDown"
		}else if win.Pressed(pixelgl.KeyUp){
			player.Direction = "LeftUp"
		}else{
			player.Direction = "Left"
		}
		player.Body.X-= 3
	}
	if win.Pressed(pixelgl.KeyRight) {
		if win.Pressed(pixelgl.KeyDown){
			player.Direction = "RightDown"
		}else if win.Pressed(pixelgl.KeyUp){
			player.Direction = "RightUp"
		}else{
			player.Direction = "Right"
		}
		player.Body.X += 3
	}
	if win.Pressed(pixelgl.KeyUp) {
		if win.Pressed(pixelgl.KeyRight){
			player.Direction = "RightUp"
		}else if win.Pressed(pixelgl.KeyLeft){
			player.Direction = "LeftUp"
		}else{
			player.Direction = "Up"
		}
		player.Body.Y += 3
	}
	if win.Pressed(pixelgl.KeyDown) {
		if win.Pressed(pixelgl.KeyRight){
			player.Direction = "RightDown"
		}else if win.Pressed(pixelgl.KeyLeft){
			player.Direction = "LeftDown"
		}else{
			player.Direction = "Down"
		}
		player.Body.Y -= 3
	}
	if win.Pressed(pixelgl.KeySpace){
		if !player.ShotFired {
			player.ShotFired = true
			player.Shot.Body.X = player.Body.X
			player.Shot.Body.Y = player.Body.Y
			player.Shot.Direction = player.Direction
			bulletPic, _ := LoadPicture("assets/bullet.png")
			player.Shot.Image.Set(bulletPic,bulletPic.Bounds())
			player.Shot.Image.Draw(win, pixel.IM.Moved(pixel.V(player.Body.X, player.Body.Y)))
			PlaySoundEffects("assets/bullet.mp3")
		}
	}
	if player.ShotFired == true {
		if player.Shot.Direction == "Left" {
			player.Shot.Body.X -= 5
		}
		if player.Shot.Direction == "Right" {
			player.Shot.Body.X += 5
		}
		if player.Shot.Direction == "Up"  {
			player.Shot.Body.Y += 5
		}
		if player.Shot.Direction == "Down"  {
			player.Shot.Body.Y -= 5
		}
		if player.Shot.Direction == "RightUp"{
			player.Shot.Body.X += 5
			player.Shot.Body.Y += 5
		}
		if player.Shot.Direction == "RightDown"{
			player.Shot.Body.X += 5
			player.Shot.Body.Y -= 5
		}
		if player.Shot.Direction == "LeftUp"{
			player.Shot.Body.X -= 5
			player.Shot.Body.Y += 5
		}
		if player.Shot.Direction == "LeftDown"{
			player.Shot.Body.X -= 5
			player.Shot.Body.Y -= 5
		}
		player.Shot.Image.Draw(win, pixel.IM.Moved(pixel.V(player.Shot.Body.X, player.Shot.Body.Y)))
	}
	if lives == 0{
		GameOver(win)
	}
	return player
}

// The UpdateLeveBadGuys takes a struct of a Badguy, Player and an integer.
// The function returns a struct of BadGuy and Player
// The function moves the badguy towards the player.
// The function also sets the bullet for the badGuy.
func UpdateLevelBadGuys(badguy BadGuy, player Player, i int) (BadGuy, Player){
	badguy.Destination = player.Body
	if IsBadGuyDead(badguy, player.Shot){
		badguy.Alive = false
		badguy.Image.Set(nil, pixel.Rect{})
		badguy.Body = pixel.Vec{}
		player.Shot.Body = pixel.Vec{}
		player.Shot.Image.Set(nil, pixel.Rect{})
		player.ShotFired = false
		player.Score += 10
		badGuys = append(badGuys[:i], badguy)
		PlaySoundEffects("assets/badGuyDead.mp3")
	}
	badguy.Count ++
	if badguy.Count % 100 == 0 {
		badguy = UpdateBadGuy(badguy)
		if badguy.Direction == "Left" {
			badguy.Body.X -= 4
		}
		if badguy.Direction == "Right" {
			badguy.Body.X += 4
		}
		if badguy.Direction == "Up"  {
			badguy.Body.Y += 4
		}
		if badguy.Direction == "Down"  {
			badguy.Body.Y -= 4
		}
	}
	return badguy, player
}

// The function takes a bullet struct as arguments and returns a boolean value.
// The function checks if the bullet of the bad guy has hit wall.
func BadGuysBulletHitWall(bullet Bullet) bool{
	for _, wall := range walls{
		if wall.Contains(bullet.Body){
			return true
		}
	}
	return false
}

// The function takes a struct of Bullet and Player.
// The function checks if the bullet of badGuy has hit player and returns a boolean value.
func BadGuyHitPlayer(bullet Bullet, player Player) bool{
	playerArea := pixel.R(player.Body.X - 13.0, player.Body.Y - 26.0, player.Body.X + 19.0, player.Body.Y + 21.0)
	if playerArea.Contains(bullet.Body){
		return true
	}
	return false
}

// The function UpdateBadGuy takes a struct of BadGuy.
// The function checks the player position and moves the badguy towards the player and returns a badguy struct.
func UpdateBadGuy(badGuy BadGuy) BadGuy {
	playerArea := pixel.R(badGuy.Destination.X - 13.0, badGuy.Destination.Y - 26.0, badGuy.Destination.X + 19.0, badGuy.Destination.Y + 21.0)
	if badGuy.Body.X >= playerArea.Min.X && badGuy.Body.X <= playerArea.Max.X{
		if badGuy.Destination.Y < badGuy.Body.Y{
			badGuy.Direction = "Down"
		}else {
			badGuy.Direction = "Up"
		}
	}else{
		if badGuy.Destination.X < badGuy.Body.X{
			badGuy.Direction = "Left"
		}else{
			badGuy.Direction = "Right"
		}
	}
	return badGuy
}


// The function IsPlayerDead takes a struct of Player as argument.
// The function checks if the player has been hit by bad guy bullet, if the player has touched wall or bad guy and return a boolean value.
func IsPlayerDead(player Player) bool{
	playerArea := pixel.R(player.Body.X - 13.0, player.Body.Y - 26.0, player.Body.X + 19.0, player.Body.Y + 21.0)
	for _, badGuy := range badGuys{
		badGuyArea := pixel.R(badGuy.Body.X - 20.0, badGuy.Body.Y - 26.0, badGuy.Body.X + 21.0, badGuy.Body.Y + 26.0)
		edgeP1 := pixel.V(player.Body.X - 13.0, player.Body.Y+21.0)
		edgeP2 := pixel.V(player.Body.X + 19.0 , player.Body.Y - 26.0)
		edgeV1 := pixel.V(badGuy.Body.X - 20.0, badGuy.Body.Y+26.0)
		edgeV2 := pixel.V(badGuy.Body.X + 21.0 , badGuy.Body.Y - 26.0)
		if playerArea.Contains(badGuyArea.Min) || playerArea.Contains(badGuyArea.Max) || playerArea.Contains(edgeV1) ||
			playerArea.Contains(edgeV2)||
			badGuyArea.Contains(playerArea.Min)|| badGuyArea.Contains(playerArea.Max) || badGuyArea.Contains(edgeP1)||
			badGuyArea.Contains(edgeP2){
			return true
		}
		for _, wall := range walls{
			if wall.Contains(playerArea.Min) || wall.Contains(playerArea.Max) || wall.Contains(edgeP1) || wall.Contains(edgeP2){
				return true
			}
		}
	}
	return false
}

// The function takes a bullet struct and a BadGuy struct as argument.
// The function checks if the bullet hit by player hits wall or bad Guy and returns a boolean value
func BulletHit(bullet Bullet, badGuy BadGuy) bool{
	badGuyArea := pixel.R(badGuy.Body.X - 20.0, badGuy.Body.Y - 26.0, badGuy.Body.X + 21.0, badGuy.Body.Y + 26.0)
	for _, wall := range walls{
		if wall.Contains(bullet.Body){
			return true
		}
	}
	if bullet.Body.Y < 0 || bullet.Body.Y > 650{
		return true
	}

	if bullet.Body.X < 0 || bullet.Body.X > 800{
		return true
	}
	if badGuyArea.Contains(bullet.Body){
		return true
	}
	return false
}

// The function takes a struct of BadGuy and Bullet as arguments.
// The function checks if bad guy has touched a wall or is hit by the player's bullet and returns a boolean value.
func IsBadGuyDead(badGuy BadGuy, bullet Bullet) bool{
	badGuyArea := pixel.R(badGuy.Body.X - 20.0, badGuy.Body.Y - 26.0, badGuy.Body.X + 21.0, badGuy.Body.Y + 26.0)
	edge1 := pixel.V(badGuy.Body.X - 20.0, badGuy.Body.Y + 26.0)
	edge2 := pixel.V(badGuy.Body.X + 21.0 , badGuy.Body.Y - 26.0)
	for _, wall := range walls{
		if wall.Contains(badGuyArea.Min) || wall.Contains(badGuyArea.Max) || wall.Contains(edge1) || wall.Contains(edge2){
			return true
		}
	}
	if badGuyArea.Contains(bullet.Body){
		return true
	}
	return false
}

// The function takes pixel window and struct of player as arguments.
// The function sets an image for the player and sets the player location to default.
func PlayerDead(win *pixelgl.Window, player Player){
	win.Clear(colornames.Dimgray)
	pic, _ := LoadPicture("assets/player_dead.png")
	player.Image.Set(pic, pic.Bounds())
	player.Image.Draw(win, pixel.IM.Moved(pixel.V(player.Body.X, player.Body.Y)))
	win.Update()
}

//The function takes pizel window and player struct as arguments
//The function will loop the number and draw image of hearys at the bottom left corner.
func DrawLives(window *pixelgl.Window, player Player){
	pic, _ := LoadPicture("assets/heart.png")
	heartSprite := pixel.NewSprite(pic,pic.Bounds())
	xloc := 30.0
	yLoc := 25.0
	for i:=0; i<player.Lives; i++ {
		heartSprite.Draw(window, pixel.IM.Moved(pixel.V(xloc,yLoc)))
		xloc += 55.0
	}
	window.Update()
}

// The function takes window as argument.
// The function will empty the window and draw game over text in the middle of the screen with the score.
func GameOver(win *pixelgl.Window){
	win.Clear(colornames.Darkolivegreen)
	walls = []pixel.Rect{}
	badGuys = []BadGuy{}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(360, 350), atlas)
	basicTxt.Color = color.Black
	fmt.Fprintln(basicTxt, "GAME OVER")
	basicTxt.Draw(win, pixel.IM)

	scoreAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreTxt := text.New(pixel.V(360, 320), scoreAtlas)
	scoreTxt.Color = color.Black
	fmt.Fprintln(scoreTxt, "SCORE:", score)
	basicTxt.Draw(win, pixel.IM)
	scoreTxt.Draw(win, pixel.IM)

}

//https://github.com/faiface/beep/blob/master/examples/playing/mp3-playing.go
// The function take and string path as argument.
// The function will play the sound clip
func PlaySoundEffects(path string){
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	s, format, _ := mp3.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(s, beep.Callback(func() {

	})))
}

// The function checks if all the bad guys are dead and returns a boolean value.
func BadGuysAllDead() bool{
	for _, badguy := range badGuys{
		if badguy.Alive == true{
			return false
		}
	}
	return true
}

// The function takes window and player struct as arguments.
// The function setsup the wall for the next level and sets the legal location for bad guys to appear.
func ChangeLevels(win *pixelgl.Window, player Player){
	player.Body = pixel.Vec{40.0, 384.0}
	win.Clear(colornames.Dimgray)
	if level == 2{
		xLoc = []float64{50, 100, 140, 180, 220, 260, 340, 380, 420, 520, 560, 600, 640, 680,720}
		yLoc = []float64{100, 140, 180, 220, 260, 300, 340, 380, 420, 460, 500, 520, 560, 600}

		SetSecondLevelWall()
		Run(win, player, GameLevel{14,13})
	}
	if level == 3{
		xLoc = []float64{40, 80, 120, 160, 200, 240, 280, 345, 385, 425, 465, 530, 570, 610, 650, 690, 735, 760}
		yLoc = []float64{100, 140, 170, 265, 300, 340, 380, 420, 460, 500, 540, 600}
		SetLevelThreeWall()
		Run(win, player, GameLevel{18, 11})
	}

	if level == 4{
		xLoc = []float64{40, 90, 130, 200, 240, 280, 345, 385, 425, 465, 530, 570, 600, 680, 720, 760}
		yLoc = []float64{100, 140, 180, 220, 260, 300, 340, 380, 420, 460, 500, 540, 600}
		SetLevelFourWall()
		Run(win, player, GameLevel{15, 12})
	}
}

//Draw walls for Level One and set up wall as rectangle
func SetLevelOneWall(){
	imd.Color = pixel.RGB(0,0,0)
	imd.Push(pixel.V(10, 60))
	imd.Push(pixel.V(325,60))
	imd.Push(pixel.V(325,70))
	imd.Push(pixel.V(10, 70))
	imd.Polygon(0)
	walls = append(walls, pixel.R(10,60,325, 70))

	imd1.Color = pixel.RGB(0,0,0)
	imd1.Push(pixel.V(475, 60))
	imd1.Push(pixel.V(790,60))
	imd1.Push(pixel.V(790,70))
	imd1.Push(pixel.V(475, 70))
	imd1.Polygon(0)
	walls = append(walls, pixel.R(475,60,790, 70))

	imd2.Color = pixel.RGB(0,0,0)
	imd2.Push(pixel.V(10, 625))
	imd2.Push(pixel.V(10,60))
	imd2.Push(pixel.V(20,60))
	imd2.Push(pixel.V(20, 625))
	imd2.Polygon(0)
	walls = append(walls, pixel.R(10,60,20, 625))

	imd3.Color = pixel.RGB(0,0,0)
	imd3.Push(pixel.V(10, 635))
	imd3.Push(pixel.V(325,635))
	imd3.Push(pixel.V(325,625))
	imd3.Push(pixel.V(10, 625))
	imd3.Polygon(0)
	walls = append(walls, pixel.R(10,625,325, 635))

	imd4.Color = pixel.RGB(0,0,0)
	imd4.Push(pixel.V(475, 635))
	imd4.Push(pixel.V(790,635))
	imd4.Push(pixel.V(790,625))
	imd4.Push(pixel.V(475, 625))
	imd4.Polygon(0)
	walls = append(walls, pixel.R(475,625,790, 635))

	imd5.Color = pixel.RGB(0,0,0)
	imd5.Push(pixel.V(780, 635))
	imd5.Push(pixel.V(780,60))
	imd5.Push(pixel.V(790,60))
	imd5.Push(pixel.V(790, 635))
	imd5.Polygon(0)
	walls = append(walls, pixel.R(780, 60,790, 635))

	imd6.Color = pixel.RGB(0,0,0)
	imd6.Push(pixel.V(190, 450))
	imd6.Push(pixel.V(190,250))
	imd6.Push(pixel.V(205, 250))
	imd6.Push(pixel.V(205,450))
	imd6.Polygon(0)
	walls = append(walls, pixel.R(190,250,205, 450))

	imd7.Color = pixel.RGB(0,0,0)
	imd7.Push(pixel.V(190, 355))
	imd7.Push(pixel.V(590,355))
	imd7.Push(pixel.V(590, 370))
	imd7.Push(pixel.V(190,370))
	imd7.Polygon(0)
	walls = append(walls, pixel.R(190,355,590, 370))

	imd8.Color = pixel.RGB(0,0,0)
	imd8.Push(pixel.V(590, 450))
	imd8.Push(pixel.V(590,250))
	imd8.Push(pixel.V(605,250))
	imd8.Push(pixel.V(605, 450))
	imd8.Polygon(0)
	walls = append(walls, pixel.R(590,250,605, 450))
}

//Draw walls for Level Two and Edit to walls from Level One
func SetSecondLevelWall(){
	walls = []pixel.Rect{}

	imd = imdraw.New(nil)
	imd.Color = pixel.RGB(0,0,0)
	imd.Push(pixel.V(10, 60))
	imd.Push(pixel.V(790,60))
	imd.Push(pixel.V(790,70))
	imd.Push(pixel.V(10, 70))
	imd.Polygon(0)
	walls = append(walls, pixel.R(10,60,790, 70))

	imd1 = imdraw.New(nil)
	imd1.Color = pixel.RGB(0,0,0)
	imd1.Push(pixel.V(10, 625))
	imd1.Push(pixel.V(10,60))
	imd1.Push(pixel.V(20,60))
	imd1.Push(pixel.V(20, 625))
	imd1.Polygon(0)
	walls = append(walls, pixel.R(10,60,20, 625))

	imd2 = imdraw.New(nil)
	imd2.Color = pixel.RGB(0,0,0)
	imd2.Push(pixel.V(10, 635))
	imd2.Push(pixel.V(10,625))
	imd2.Push(pixel.V(300,625))
	imd2.Push(pixel.V(300, 635))
	imd2.Polygon(0)
	walls = append(walls, pixel.R(10,625,300, 635))

	imd3 = imdraw.New(nil)
	imd3.Color = pixel.RGB(0,0,0)
	imd3.Push(pixel.V(500, 635))
	imd3.Push(pixel.V(790,635))
	imd3.Push(pixel.V(790,625))
	imd3.Push(pixel.V(500, 625))
	imd3.Polygon(0)
	walls = append(walls, pixel.R(500,625,790, 635))

	imd4 = imdraw.New(nil)
	imd4.Color = pixel.RGB(0,0,0)
	imd4.Push(pixel.V(780, 635))
	imd4.Push(pixel.V(780,60))
	imd4.Push(pixel.V(790,60))
	imd4.Push(pixel.V(790, 635))
	imd4.Polygon(0)
	walls = append(walls, pixel.R(780, 60,790, 635))

	imd5 = imdraw.New(nil)
	imd5.Color = pixel.RGB(0,0,0)
	imd5.Push(pixel.V(300, 635))
	imd5.Push(pixel.V(315,635))
	imd5.Push(pixel.V(315,430))
	imd5.Push(pixel.V(300, 430))
	imd5.Polygon(0)
	walls = append(walls, pixel.R(300, 430,315, 635))

	imd6 = imdraw.New(nil)
	imd6.Color = pixel.RGB(0,0,0)
	imd6.Push(pixel.V(485, 635))
	imd6.Push(pixel.V(500,635))
	imd6.Push(pixel.V(500, 430))
	imd6.Push(pixel.V(485,430))
	imd6.Polygon(0)
	walls = append(walls, pixel.R(485,430,500, 635))

	imd7 = imdraw.New(nil)
	imd7.Color = pixel.RGB(0,0,0)
	imd7.Push(pixel.V(300, 265))
	imd7.Push(pixel.V(315,265))
	imd7.Push(pixel.V(315, 60))
	imd7.Push(pixel.V(300,60))
	imd7.Polygon(0)
	walls = append(walls, pixel.R(300,60,315, 265))

	imd8 = imdraw.New(nil)
	imd8.Color = pixel.RGB(0,0,0)
	imd8.Push(pixel.V(485, 265))
	imd8.Push(pixel.V(500,265))
	imd8.Push(pixel.V(500,60))
	imd8.Push(pixel.V(485, 60))
	imd8.Polygon(0)
	walls = append(walls, pixel.R(485,60,500, 265))
}

//Draw walls for Level Three and Edit to walls from Level Two
func SetLevelThreeWall(){
	walls = []pixel.Rect{}
	walls = append(walls, pixel.R(10,60,790, 70))
	walls = append(walls, pixel.R(10,60,20, 625))
	walls = append(walls, pixel.R(10,625,300, 635))
	walls = append(walls, pixel.R(500,625,790, 635))
	walls = append(walls, pixel.R(780, 60,790, 635))
	walls = append(walls, pixel.R(300, 430,315, 635))
	walls = append(walls, pixel.R(485,430,500, 635))

	imd7 = imdraw.New(nil)
	imd7.Color = pixel.RGB(0,0,0)
	imd7.Push(pixel.V(10,210))
	imd7.Push(pixel.V(10,225))
	imd7.Push(pixel.V(315, 225))
	imd7.Push(pixel.V(315,210))
	imd7.Polygon(0)
	walls = append(walls, pixel.R(10,210,315, 225))

	imd8 = imdraw.New(nil)
	imd8.Color = pixel.RGB(0,0,0)
	imd8.Push(pixel.V(485, 225))
	imd8.Push(pixel.V(485,210))
	imd8.Push(pixel.V(780,210))
	imd8.Push(pixel.V(780, 225))
	imd8.Polygon(0)
	walls = append(walls, pixel.R(485,210,780, 225))
}

//Draw walls for Level Four and Edit to walls from Level Three
func SetLevelFourWall(){
	walls = []pixel.Rect{}

	walls = append(walls, pixel.R(10,60,790, 70))

	imd1 = imdraw.New(nil)
	imd1.Color = pixel.RGB(0,0,0)
	imd1.Push(pixel.V(10, 250))
	imd1.Push(pixel.V(10,60))
	imd1.Push(pixel.V(20,60))
	imd1.Push(pixel.V(20, 250))
	imd1.Polygon(0)
	walls = append(walls, pixel.R(10,60,20, 250))

	imd2 = imdraw.New(nil)
	imd2.Color = pixel.RGB(0,0,0)
	imd2.Push(pixel.V(10, 400))
	imd2.Push(pixel.V(10,635))
	imd2.Push(pixel.V(20,635))
	imd2.Push(pixel.V(20, 400))
	imd2.Polygon(0)
	walls = append(walls, pixel.R(10,400,20, 635))

	imd3 = imdraw.New(nil)
	imd3.Color = pixel.RGB(0,0,0)
	imd3.Push(pixel.V(10, 625))
	imd3.Push(pixel.V(10,635))
	imd3.Push(pixel.V(790,635))
	imd3.Push(pixel.V(790, 625))
	imd3.Polygon(0)
	walls = append(walls, pixel.R(10,625,790, 635))

	imd4 = imdraw.New(nil)
	imd4.Color = pixel.RGB(0,0,0)
	imd4.Push(pixel.V(780, 400))
	imd4.Push(pixel.V(790,400))
	imd4.Push(pixel.V(790,635))
	imd4.Push(pixel.V(780, 635))
	imd4.Polygon(0)
	walls = append(walls, pixel.R(780, 400,790, 635))

	imd5 = imdraw.New(nil)
	imd5.Color = pixel.RGB(0,0,0)
	imd5.Push(pixel.V(780, 60))
	imd5.Push(pixel.V(790,60))
	imd5.Push(pixel.V(790,250))
	imd5.Push(pixel.V(780, 250))
	imd5.Polygon(0)
	walls = append(walls, pixel.R(780, 60,790, 250))

	imd6 = imdraw.New(nil)
	imd6.Color = pixel.RGB(0,0,0)
	imd6.Push(pixel.V(625, 160))
	imd6.Push(pixel.V(640,160))
	imd6.Push(pixel.V(640, 525))
	imd6.Push(pixel.V(625,525))
	imd6.Polygon(0)
	walls = append(walls, pixel.R(625,160,640, 525))

	imd7 = imdraw.New(nil)
	imd7.Color = pixel.RGB(0,0,0)
	imd7.Push(pixel.V(160, 160))
	imd7.Push(pixel.V(175,160))
	imd7.Push(pixel.V(175, 525))
	imd7.Push(pixel.V(160,525))
	imd7.Polygon(0)
	walls = append(walls, pixel.R(160,160,175, 525))

	imd8 = imdraw.New(nil)
}

//Set walls and Bad Guys to empty slice and draw window that says "Game Over" and display Total Score.
func GameWon(win *pixelgl.Window, player Player){
	win.Clear(colornames.Cadetblue)
	walls = []pixel.Rect{}
	badGuys = []BadGuy{}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(360, 350), atlas)
	basicTxt.Color = color.Black
	fmt.Fprintln(basicTxt, "YOU WON !!!")
	scoreAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreTxt := text.New(pixel.V(360, 320), scoreAtlas)
	scoreTxt.Color = color.Black
	fmt.Fprintln(scoreTxt, "SCORE:", player.Score)
	basicTxt.Draw(win, pixel.IM)
	scoreTxt.Draw(win, pixel.IM)
}