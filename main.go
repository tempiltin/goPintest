package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Fyne App yaratish
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())
	myWindow := myApp.NewWindow("Web Pentesting Tool")
	myWindow.Resize(fyne.NewSize(600, 400))

	// Input maydoni
	input := widget.NewEntry()
	input.SetPlaceHolder("Web sayt havolasini kiriting...")

	// Chiqish maydoni
	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Natijalar bu yerda ko'rinadi...")
	output.Wrapping = fyne.TextWrapWord

	// "Tahlil qilish" tugmasi
	button := widget.NewButton("Tahlil qilish", func() {
		site := input.Text
		if site == "" {
			output.SetText("Havola kiritilmadi!")
			return
		}

		// Tahlil natijalarini olish
		result := runPentest(site)
		output.SetText(result)
	})

	// GUI joylashuvi
	content := container.NewVBox(
		widget.NewLabel("Web Pentesting Tool"),
		input,
		button,
		output,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func runPentest(site string) string {
	var result strings.Builder
	result.WriteString("Web sayt: " + site + "\n\n")

	// IP manzilni aniqlash
	ips, err := net.LookupIP(site)
	if err != nil {
		result.WriteString("IP manzilni topib bo'lmadi!\n")
	} else {
		for _, ip := range ips {
			result.WriteString("IP: " + ip.String() + "\n")
		}
	}

	// Sub-domenlarni aniqlash (os/exec orqali `subfinder` kabi vositadan foydalanish mumkin)
	subdomains, err := exec.Command("subfinder", "-d", site).Output()
	if err != nil {
		result.WriteString("Sub-domenlarni topishda xatolik yuz berdi!\n")
	} else {
		result.WriteString("Sub-domenlar:\n" + string(subdomains) + "\n")
	}

	// Portlarni skanerlash (80 va 443 portlari uchun)
	ports := []int{80, 443}
	result.WriteString("Ochiq portlar:\n")
	for _, port := range ports {
		address := fmt.Sprintf("%s:%d", site, port)
		conn, err := net.Dial("tcp", address)
		if err == nil {
			result.WriteString(fmt.Sprintf("Port %d: Ochiq\n", port))
			conn.Close()
		} else {
			result.WriteString(fmt.Sprintf("Port %d: Yopiq\n", port))
		}
	}

	return result.String()
}
