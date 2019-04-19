/*
Saurav P. Shrestha
00369895
Final Project
*/
package main

import (
	"github.com/faiface/pixel"
	"testing"
)

func TestNewPlayer(t *testing.T){
	if NewPlayer().Direction != "Right"{
		t.Error("Test Failed")
	}
	if NewPlayer().ShotFired != false{
		t.Error("Test failed")
	}
	if NewPlayer().Lives != 3{
		t.Error("Test failed")
	}
	if NewPlayer().Score != 0{
		t.Error("Test failed")
	}
}

func TestNewBadGuy(t *testing.T){
	if NewBadGuy().Alive == false{
		t.Error("Test failed")
	}
	if NewBadGuy().ShotFired == true{
		t.Error("Test failed")
	}
}

func TestNewBullet(t *testing.T){
	if NewBullet().Direction != ""{
		t.Error("Test failed")
	}
	if NewBullet().Hit == true{
		t.Error("Test failed")
	}
}

func TestLoadPicture(t *testing.T) {
	_ , err := LoadPicture("assets/error.png")
	if err.Error() != "open assets/error.png: The system cannot find the file specified."{
		t.Error("Test failed")
	}

	pic , err1 := LoadPicture("assets/bullet.png")
	if err1 != nil {
		t.Error("Test failed")
	}
	if pic.Bounds().Max.X != 22{
		t.Error("Test Failed")
	}
}

func TestUpdateLevelBadGuys(t *testing.T) {
	player := NewPlayer()
	badGuy := NewBadGuy()

	bg, pl := UpdateLevelBadGuys(badGuy, player, 0)

	if bg.Destination.X != pl.Body.X{
		t.Error("Test failed")
	}
	if bg.Destination.Y != pl.Body.Y{
		t.Error("Test failed")
	}
}

func TestBadGuysBulletHitWall(t *testing.T) {
	walls = append(walls, pixel.R(10,60,325, 70))
	bullet := NewBullet()
	bullet.Body = pixel.Vec{10,65}
	if BadGuysBulletHitWall(bullet) != true {
		t.Error("Test failed")
	}
	bullet.Body = pixel.Vec{200,200}
	if BadGuysBulletHitWall(bullet) != false{
		t.Error("Test failed")
	}
}

func TestBadGuyHitPlayer(t *testing.T) {
	player := NewPlayer()
	bullet := NewBullet()
	bullet.Body = pixel.Vec{50, 390}
	if BadGuyHitPlayer(bullet, player) != true{
		t.Error("Test failed")
	}
	bullet.Body = pixel.Vec{100,200}
	if BadGuyHitPlayer(bullet, player) != false{
		t.Error("Test failed")
	}
}

func TestUpdateBadGuy(t *testing.T) {
	badGuy := NewBadGuy()
	badGuy.Body = pixel.Vec{200, 400}
	badGuy.Destination = pixel.Vec{100, 200}
	if UpdateBadGuy(badGuy).Direction != "Left"{
		t.Error("Test failed")
	}
	badGuy.Body = pixel.Vec{100, 400}
	badGuy.Destination = pixel.Vec{100, 200}
	if UpdateBadGuy(badGuy).Direction != "Down"{
		t.Error("Test failed")
	}
}

func TestIsPlayerDead(t *testing.T) {
	walls = append(walls, pixel.R(10,60,325, 70))
	player := NewPlayer()

	if IsPlayerDead(player) != false{
		t.Error("Test failed")
	}
	player.Body = pixel.Vec{20, 40}
	if IsPlayerDead(player) != true{
		t.Error("Test failed")
	}
}

func TestBulletHit(t *testing.T) {
	bullet := NewBullet()
	bullet.Body = pixel.Vec{40,50}
	badGuy := NewBadGuy()
	walls = append(walls, pixel.R(10,60,325, 70))
	if BulletHit(bullet, badGuy) != false{
		t.Error("Test failed")
	}
	bullet.Body = pixel.Vec{10,65}
	if BulletHit(bullet, badGuy) != true{
		t.Error("Test failed")
	}
}

func TestIsBadGuyDead(t *testing.T) {
	badGuy := NewBadGuy()
	bullet := NewBullet()
	walls = append(walls, pixel.R(10,60,325, 70))
	badGuy.Body = pixel.Vec{150, 250}
	if IsBadGuyDead(badGuy, bullet) != false{
		t.Error("Test failed")
	}
	bullet.Body = pixel.Vec{140, 240}
	if IsBadGuyDead(badGuy, bullet) != true{
		t.Error("Test failed")
	}
}

func TestBadGuysAllDead(t *testing.T) {
	badGuy := NewBadGuy()
	badGuys = append(badGuys, badGuy)
	if BadGuysAllDead() != false{
		t.Error("Test failed")
	}

	badGuys = []BadGuy{}
	badGuy.Alive = false
	badGuys = append(badGuys, badGuy)
	if BadGuysAllDead() != true{
		t.Error("Test failed")
	}
}