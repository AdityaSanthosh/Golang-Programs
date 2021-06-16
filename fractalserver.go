package main

import (
	"html/template"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
	"os"
	"strconv"
)

func fractalserver(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	xmin, _ := strconv.Atoi( r.PostFormValue("xmin")[0:])
	ymin, _ := strconv.Atoi( r.PostFormValue("ymin")[0:])
	xmax, _ := strconv.Atoi( r.PostFormValue("xmax")[0:])
	ymax, _ := strconv.Atoi( r.PostFormValue("ymax")[0:])
	const (
		width, height = 1024,1024
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	out, _:=os.Create("mandelbrot.png")
	for py := 0;py<height;py++ {
		y := float64(py)/float64(height)*float64(ymax-ymin)+ float64(ymin)
		for px:=0;px < width;px++ {
			x := float64(px)/float64(width)*float64(xmax-xmin)+ float64(xmin)
			z:= complex(x,y)
			img.Set(px,py,mandelbrot(z))
		}
	}
	png.Encode(out, img)
	out.Close()
	image, err := os.Open("mandelbrot.png")
	if err != nil {
		log.Fatal(err) // perhaps handle this nicer
	}
	defer image.Close()
	w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(w, image)
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n:= uint8(0);n<iterations;n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255-contrast*n}
		}
	}
	return color.Black
}

func homepage(w http.ResponseWriter, r *http.Request) {
	files, err := template.ParseFiles("fractalclient.html")
	if err != nil {
		return
	}
	err=files.Execute(w,"fractalclient.html")
}

func main() {
	http.HandleFunc("/",homepage)
	http.HandleFunc("/render/",fractalserver)
	http.ListenAndServe(":8000",nil)
}