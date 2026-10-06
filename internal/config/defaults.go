package config

// DefaultConfig returns the default configuration mirroring smford.github.io.
func DefaultConfig() Config {
	return Config{
		Operator: OperatorConfig{
			Handle:   "smford",
			Host:     "mainframe",
			Title:    "smford.github.io",
			Subtitle: "> WOODWORKING SOFTWARE ENGINEER // HARDWARE PROTOTYPER",
			Bio:      "Bridging low-level firmware engineering, bespoke IoT hardware controllers, real-time laser telemetry, and physical CNC/timber fabrication from the Bay Area network.",
			SecNet:   "ACTIVE",
		},
		Telemetry: TelemetryConfig{
			Kernel:       "6.1.0-DEBIAN",
			Cyberware:    "ESP32 / PI / CAD",
			SecurityICE:  "UNRESTRICTED",
			Repositories: "27 DEPLOYED",
		},
		Settings: SettingsConfig{
			DefaultTheme: "cyberpunk",
			MouseEnabled: true,
			SkipBoot:     false,
		},
		VFS: map[string]string{
			"bio.txt": `OPERATOR PROFILE:
Name: Steve Ford (@smford)
Role: Woodworking Software Engineer & Hardware Prototyper
Specialty: Bridging physical timber engineering with embedded C++, Linux sysadmin, and custom CAD tooling.
Location: Bay Area Node // UK Matrix Grid`,

			"hardware.ino": `// EEH ESP32 RFID Safety Sentry
#include <SPI.h>
#include <MFRC522.h>

#define SS_PIN    5
#define RST_PIN   22
#define RELAY_PIN 25

MFRC522 rfid(SS_PIN, RST_PIN);

void setup() {
  Serial.begin(115200);
  SPI.begin();
  rfid.PCD_Init();
  pinMode(RELAY_PIN, OUTPUT);
  digitalWrite(RELAY_PIN, HIGH); // Lockout active
  Serial.println("[MAINFRAME] Safety interlock armed.");
}

void loop() {
  if (!rfid.PICC_IsNewCardPresent() || !rfid.PICC_ReadCardSerial()) return;
  Serial.println("[AUTH] Clearance granted. Engaging power relay.");
  digitalWrite(RELAY_PIN, LOW);
  delay(10000); // 10s tool duty cycle
  digitalWrite(RELAY_PIN, HIGH);
}`,

			"temp_monitor.ino": `// ESP32 Temperature Regulator & Alert Daemon
#include <WiFi.h>
#include <OneWire.h>
#include <DallasTemperature.h>

#define ONE_WIRE_BUS 4
#define FRIDGE_RELAY 26

OneWire oneWire(ONE_WIRE_BUS);
DallasTemperature sensors(&oneWire);

void setup() {
  pinMode(FRIDGE_RELAY, OUTPUT);
  sensors.begin();
}

void loop() {
  sensors.requestTemperatures();
  float tempC = sensors.getTempCByIndex(0);
  if (tempC > 5.0) {
    digitalWrite(FRIDGE_RELAY, LOW); // Active cooling
  } else if (tempC < 3.5) {
    digitalWrite(FRIDGE_RELAY, HIGH); // Disengage cooling
  }
  delay(5000);
}`,

			"snmpd.conf": `# Raspberry Pi SNMPD Telemetry Configuration
syslocation "Workshop Bay 4"
syscontact "smford <admin@mainframe.local>"
rocommunity public default
# Hardware sensors script
extend cpu_temp /usr/bin/vcgencmd measure_temp`,

			"notes.md": `# Night City Hardware Bench
- [x] Tested Bambu magnetic poop chute remix on X1-Carbon
- [x] Verified async upload progress bar on ESP32 SPIFFS
- [x] Aligned wide-angle optics on Pi Zero RTSP laser-cam
- [x] Calibrated Chandra-X black hole raymarching shader (smford/blackhole)
- [x] Optimized iTerm2 OSC 1337 image rendering in mdee (smford/mdee)
- [x] Deployed YAML anchor validator and RBAC manager (smford/yaml-web-validator)
- [x] Verified high-res Sixel & iTerm2 terminal Mermaid diagrams (smford/golang-mermaid)
- [x] Published official Homebrew tap formulas (smford/homebrew-tap)
- [x] Deployed IPv4/IPv6 CIDR decomposition & collision engine (smford/cidr-calculator)
- [x] 3D printed emergency stop wall mount (smford/emergency-stop-holder)
- [ ] Install Pica Dry wall holsters on main workbench`,
		},
		QuickChips: []QuickChip{
			{Label: "help", Command: "help"},
			{Label: "projects", Command: "projects"},
			{Label: "matrix-rain", Command: "matrix-rain"},
			{Label: "cidr", Command: "cidr 192.168.1.10+5"},
			{Label: "blackhole", Command: "blackhole"},
			{Label: "moon", Command: "moon"},
			{Label: "bbs", Command: "bbs"},
			{Label: "burger", Command: "burger"},
			{Label: "mdee", Command: "mdee notes.md"},
			{Label: "mcat", Command: "mcat notes.md"},
			{Label: "delim", Command: "delim"},
			{Label: "toolsplit", Command: "toolsplit calc 199 99"},
			{Label: "rbac audit", Command: "rbac audit"},
			{Label: "yaml lint", Command: "yaml lint"},
			{Label: "license mit", Command: "license mit"},
			{Label: "brew list", Command: "brew list"},
			{Label: "mermaid", Command: "mermaid"},
			{Label: "snmp", Command: "snmp"},
			{Label: "rfid", Command: "rfid"},
			{Label: "temp", Command: "temp"},
			{Label: "cam", Command: "cam"},
			{Label: "cad", Command: "cad"},
			{Label: "estop", Command: "estop"},
			{Label: "cd firmware", Command: "cd firmware"},
			{Label: "ls", Command: "ls"},
			{Label: "clear", Command: "clear"},
		},
		CustomCommands: map[string]Command{
			"blackhole": {
				Description: "Relativistic Schwarzschild black hole raymarcher & ISCO telemetry",
				Output: `[CHANDRA-X ACCRETION SIMULATOR]
Coordinate Observer: Schwarzschild Metric (G=c=1)
Event Horizon (r_s): 29.53 km | ISCO: 88.59 km
Photon Sphere: 44.30 km | Gravitational Redshift: z=0.414
Launch simulation: use 'open blackhole' or 'blackhole web'`,
				URL: "https://stephenford.org/blackhole/",
			},
			"moon": {
				Description: "3D lunar reconnaissance ephemeris & Apollo sites",
				Output: `[SELENOGRAPHIC RECONNAISSANCE SUITE]
Target: Lunar Surface Reconnaissance (LOD: Level 14 Mesh)
Sub-Solar Point: 0.84°N 14.22°E | Optical Libration: 5.12°
Telemetry stream: use 'open moon' or 'moon web'`,
				URL: "https://stephenford.org/fly-around-the-moon/",
			},
			"bbs": {
				Description: "Dial into Clark Dev PCBoard v15.22 MS-DOS BBS with audio",
				Output: `[PCBOARD v15.22 MODEM HANDSHAKE]
Connecting to 28,800 bps V.34 modem pool...
CONNECT 28800 / V.42bis / CP437 ANSI
Welcome to The Byte Bunker BBS Node #1 (Sysop: smford)
Launch interface: use 'open bbs' or 'bbs web'`,
				URL: "https://stephenford.org/pcboard-web-emu/",
			},
			"burger": {
				Description: "3D culinary exploded view & layer inspection",
				Output: `[CULINARY GASTRONOMY EXPLODED SCHEMATIC]
Layers parsed: Brioche Crown, Microgreens, Aged Cheddar, Smashed Patty, Aioli, Brioche Heel.
Launch interface: use 'open burger' or 'burger web'`,
				URL: "https://stephenford.org/hamburger-decomp/",
			},
			"toolsplit": {
				Description: "Power tool promotional discount calculator and prorated return receipt engine",
				Output: `[TOOLSPLIT BOGO SAVINGS OPTIMIZER]
Tool A (Primary): $199.00 | Tool B (Free): $99.00
Combined MSRP: $298.00 | Promotion Bundle Cost: $199.00
Effective Discount: 33.22%
Adjusted Return Value A: $132.89 | Adjusted Return Value B: $66.11
Launch interactive calculator: use 'open toolsplit'`,
				URL: "https://stephenford.org/toolsplit/",
			},
			"snmp": {
				Description: "Raspberry Pi CPU temperature & system telemetry MIBs",
				Output: `[HOST-RESOURCES-MIB TELEMETRY]
HOST-RESOURCES-MIB::hrSystemUptime.0 = (439281200) 50 days, 20:13:32.00
UCD-SNMP-MIB::extNames.1 = STRING: "cpu_temp"
UCD-SNMP-MIB::extOutput.1 = STRING: "temp=42.8'C"
NET-SNMP-EXTEND-MIB::nsExtendOutputFull."cpu_temp" = STRING: "42.8 C"`,
			},
			"esp32": {
				Description: "ESP32 AsyncWebServer chunked SPIFFS upload telemetry",
				Output: `[ESP32 ASYNC WEBSERVER FILE-UPLOAD ENGINE]
Heap Free: 218,440 bytes | SPIFFS Total: 1,441,792 bytes (Used: 382,100 bytes)
Async TCP Handshake: 120ms | WebSocket Clients Active: 2
Status: READY FOR MULTI-PART CHUNKED STREAMS`,
				URL: "https://github.com/smford/esp32-asyncwebserver-fileupload-example",
			},
			"cam": {
				Description: "Pi Zero W RTSP low-latency CNC surveillance camera",
				Output: `[CNC OPTICAL LASER-CAM PIPELINE]
Stream: rtsp://192.168.4.180:8554/unicast
Encoder: Hardware H.264 V4L2 Device (/dev/video0)
Resolution: 1920x1080 @ 25fps (Sub-120ms Latency)
Hardware Invert: FLIP_H=1 FLIP_V=1 (Ceiling Gantry Mount)`,
				URL: "https://github.com/smford/raspberry-pi-streaming-camera",
			},
			"rfid": {
				Description: "East Essex Hackspace 3kW machinery safety sentry",
				Output: `[EAST ESSEX HACKSPACE INTERLOCK SENTRY]
Bus: SPI 13.56MHz (RC522 Module on GPIO_05/18/23)
Contactor Isolation: PC817 Optocoupler -> 30A AC Relay
Authorized Badges: 48 Registered Operators
Status: INTERLOCK LOCKED (Awaiting induction badge swipe)`,
				URL: "https://github.com/smford/eeh-esp32-rfid",
			},
			"temp": {
				Description: "DS18B20 tortoise hibernation regulation daemon",
				Output: `[THERMAL REGULATION TELEMETRY]
Probe 1 (Chamber Center): 4.12°C [OPTIMAL 3.5°C - 5.0°C]
Probe 2 (Substrate):      3.98°C [OPTIMAL]
Relay SSR (GPIO 26):      OFF (Passive Thermal Equilibrium)
Telemetry Sink: InfluxDB / Grafana (Push interval: 30s)`,
				URL: "https://github.com/smford/esp32-temp-monitor",
			},
			"cad": {
				Description: "Bambu 3D CAD prints & parametric enclosures",
				Output: `[PARAMETRIC CAD FABRICATION REPOSITORY]
Chassis 1: Sunton-ESP32-2432S028R-Case (Bambu PLA-CF, M3 Heat-Set)
Chassis 2: Bambu Lab Purge Chute Magnetic Snap-Lock (Neodymium N52)
Chassis 3: 22mm Industrial E-Stop Wall Bracket (75mm & 92mm specs)`,
				URL: "https://github.com/smford/Sunton-ESP32-2432S028R-Case",
			},
			"estop": {
				Description: "Workshop emergency stop mount & safety housing",
				Output: `[EMERGENCY STOP SAFETY SYSTEM]
Actuator: 22mm Turn-to-Reset Mushroom Head Pushbutton
Mounting: Bambu CF-PETG Rigid Flange with Countersunk M3 Holes
Wiring: Dual NC contacts in series with contactor coil circuit`,
				URL: "https://github.com/smford/emergency-stop-holder",
			},
		},
		Skills: []Skill{
			{Name: "EMBEDDED C++ / ESP32", Percentage: 96},
			{Name: "LINUX / DEBIAN SYSADMIN", Percentage: 94},
			{Name: "CAD & CNC TIMBER FABRICATION", Percentage: 92},
			{Name: "GOLANG & NETWORKING", Percentage: 88},
		},
		Benches: []Bench{
			{
				ID:    "01",
				Title: "RAPID FABRICATION",
				Items: []string{
					"Bambu Lab 3D Printing (CF-Reinforced PETG/PLA)",
					"Custom CNC Gantry Router & Laser Cutter",
					"Parametric CAD (Fusion 360 / OpenSCAD)",
				},
			},
			{
				ID:    "02",
				Title: "EMBEDDED LAB",
				Items: []string{
					"ESP32 Dual-Core (Xtensa 240MHz, SPIFFS/LittleFS)",
					"Raspberry Pi Zero W / Compute Module 4",
					"Digital Logic Analyzers & Articulating Inspection",
				},
			},
			{
				ID:    "03",
				Title: "NETWORK & TELEMETRY",
				Items: []string{
					"Ubiquiti UniFi Local PoE Network Backbone",
					"MikroTik RouterOS Telemetry Hooks",
					"Home Assistant & Zigbee Mesh Integration",
				},
			},
		},
		BuildLogs: []BuildLog{
			{
				ID:       "log-042",
				Category: "Hardware",
				Title:    "EEH ESP32 RFID Tool Sentry",
				Date:     "2026-04-12",
				Status:   "DEPLOYED_ACTIVE",
				Summary:  "Industrial RFID safety interlock system preventing unauthorized activation of 3kW workshop machinery and high-draw fabrication tools.",
				BOM: []BOMItem{
					{Item: "ESP32 DevKit V1 (38-pin)", Qty: "1", Role: "Main Microcontroller"},
					{Item: "RC522 13.56MHz RFID Module", Qty: "1", Role: "User Badge Scanner"},
					{Item: "Optocoupler 5V Relay Block", Qty: "1", Role: "High-Voltage Coil Isolation"},
					{Item: "128x64 I2C OLED Display", Qty: "1", Role: "Telemetry & Operator Prompt"},
				},
				Pinout: []PinoutItem{
					{GPIO: "GPIO 21", Device: "I2C SDA", Note: "OLED Screen Data line"},
					{GPIO: "GPIO 22", Device: "I2C SCL", Note: "OLED Screen Clock line"},
					{GPIO: "GPIO 5", Device: "SPI SS / CS", Note: "RC522 RFID Chip Select"},
					{GPIO: "GPIO 18", Device: "SPI SCK", Note: "RC522 SPI Clock"},
					{GPIO: "GPIO 25", Device: "Relay Trigger", Note: "Active-Low Contactor Fire"},
				},
				Schematic: `[ 230V AC Mains ] ---> [ 30A Contactor ] ---> [ Workshop Table Saw ]
                             ^
                             | (Galvanic Isolation)
[ ESP32 GPIO 25 ] ---> [ PC817 Optocoupler ] ---> [ 5V Relay Coil ]`,
				Notes: "To prevent inductive back-EMF surges when 3kW induction motors spin down, the 5V relay coil is completely isolated with an optocoupler and flyback diode. User badges are verified locally via flash lookup tables.",
			},
			{
				ID:       "log-089",
				Category: "Optics",
				Title:    "Pi Zero RTSP Laser-Cam",
				Date:     "2026-05-20",
				Status:   "DEPLOYED_ACTIVE",
				Summary:  "Ultra-compact overhead optical monitoring system streaming low-latency video from inside a high-power CNC laser enclosure.",
				BOM: []BOMItem{
					{Item: "Raspberry Pi Zero W", Qty: "1", Role: "Video Streaming Host"},
					{Item: "5MP Wide-Angle Cam (160° FoV)", Qty: "1", Role: "Optical Sensor"},
					{Item: "3D Printed Magnetic Enclosure", Qty: "1", Role: "Chassis & Smoke Baffle"},
				},
				Pinout: []PinoutItem{
					{GPIO: "CSI Ribbon", Device: "Camera Port", Note: "Direct 2-lane MIPI CSI-2"},
					{GPIO: "MicroUSB", Device: "Power", Note: "5V 2.5A Regulated Supply"},
				},
				Schematic: `[ Wide-Angle Cam ] -(CSI)-> [ Pi Zero W ] -(H.264)-> [ v4l2rtspserver ] -(RTSP/WiFi)-> [ Mainframe ]`,
				Notes:     "Configured v4l2rtspserver with hardware H.264 encoder flags (-H 1080 -W 1920 -F 25) achieving sub-120ms latency across local 5GHz Wi-Fi mesh.",
			},
			{
				ID:       "log-104",
				Category: "Fabrication",
				Title:    "Sunton ESP32 CAD Enclosure",
				Date:     "2026-06-18",
				Status:   "PRODUCTION_RELEASE",
				Summary:  "Parametric 3D-printed enclosure designed for the Sunton ESP32-2432S028R integrated LCD and touchscreen development deck.",
				BOM: []BOMItem{
					{Item: "Sunton ESP32-2432S028R", Qty: "1", Role: "Integrated Touch Deck"},
					{Item: "Bambu PLA-CF (Carbon Fiber)", Qty: "45g", Role: "Structural Case Shell"},
					{Item: "M3x4mm Brass Heat-Set Inserts", Qty: "4", Role: "Threaded Fasteners"},
				},
				Pinout: []PinoutItem{
					{GPIO: "Internal", Device: "ST7789 TFT", Note: "Parallel 8-bit SPI display"},
					{GPIO: "GPIO 35", Device: "LDR Sensor", Note: "Ambient Light Sensor"},
				},
				Schematic: `[ Top Bezel w/ Screen Cutout ]
          ↓ (M3 Screws)
[ Sunton ESP32 Deck ]
          ↓ (Heat-set Inserts)
[ Bottom Shell w/ VESA & DIN Rail Mounts ]`,
				Notes: "Modeled in Fusion 360 with 0.2mm tolerance offsets for snap-fit tabs. Integrated ventilation channels direct convection airflow across the onboard ESP32 Wi-Fi module.",
			},
			{
				ID:       "log-112",
				Category: "Hardware",
				Title:    "ESP32 Precision Temp Regulator",
				Date:     "2026-07-05",
				Status:   "DEPLOYED_ACTIVE",
				Summary:  "Closed-loop micro-climate regulation unit logging DS18B20 telemetry to TICK stack and cycling 230V cooling relays for precision tortoise hibernation.",
				BOM: []BOMItem{
					{Item: "ESP32 Dev Board", Qty: "1", Role: "Controller & Wi-Fi Node"},
					{Item: "DS18B20 Probes (Waterproof)", Qty: "2", Role: "Multi-point Thermal Probes"},
					{Item: "Solid State Relay 10A", Qty: "1", Role: "Compressor Switching"},
				},
				Pinout: []PinoutItem{
					{GPIO: "GPIO 4", Device: "OneWire Bus", Note: "Dallas DS18B20 Sensors"},
					{GPIO: "GPIO 26", Device: "Relay SSR", Note: "Cooling Compressor Control"},
				},
				Schematic: `[ DS18B20 Probes ] ---> [ ESP32 GPIO 4 ] ---> [ SSR Relay ] ---> [ Mini Fridge ]`,
				Notes:     "Includes built-in webserver for live threshold configuration and Pushover emergency pushes if temperatures drift +/- 1.5°C outside the safety window.",
			},
			{
				ID:       "log-118",
				Category: "Fabrication",
				Title:    "Timber Rafter & Outbuilding CAD Engine",
				Date:     "2026-07-28",
				Status:   "COMPLETED",
				Summary:  "Interactive mathematical computation suite calculating complex dual-pitch birdsmouth cuts, structural joist spans, and cladding schedules.",
				BOM: []BOMItem{
					{Item: "WebGL / HTML5 Canvas", Qty: "1", Role: "Visual Geometry Renderer"},
					{Item: "UK Building Regs Part A/H", Qty: "1", Role: "Structural Compliance Matrix"},
				},
				Pinout: []PinoutItem{
					{GPIO: "Virtual", Device: "DOM Engine", Note: "Real-time trigonometric rafter rendering"},
				},
				Schematic: `[ Span & Pitch Inputs ] ---> [ Trigonometric Solver ] ---> [ Canvas SVG Cutlist & BOM ]`,
				Notes:     "Executes sub-millimeter timber cut calculations for complex hip rafters, calculating structural dead loads and Part H drainage runoffs instantly.",
			},
		},
		Repositories: []Repository{
			{
				Name:     "esp32-asyncwebserver-fileupload-example",
				Category: "Embedded",
				Stars:    104,
				Forks:    30,
				Desc:     "High-performance asynchronous web server for ESP32 demonstrating drag-and-drop SPIFFS filesystem uploads with real-time browser progress bar telemetry.",
				Tech:     []string{"C++", "ESP32", "SPIFFS", "AsyncWebServer", "WebSockets"},
				URL:      "https://github.com/smford/esp32-asyncwebserver-fileupload-example",
			},
			{
				Name:     "raspberry-pi-streaming-camera",
				Category: "RaspberryPi",
				Stars:    83,
				Forks:    5,
				Desc:     "Ultra-low latency RTSP camera feed on Raspberry Pi Zero W with hardware video flip (v4l2rtspserver + Nginx) for real-time laser cutter enclosure monitoring.",
				Tech:     []string{"Raspbian", "v4l2rtspserver", "Nginx", "Laser CNC", "MIPI-CSI"},
				URL:      "https://github.com/smford/raspberry-pi-streaming-camera",
			},
			{
				Name:     "rpi-snmpd-configuration",
				Category: "RaspberryPi",
				Stars:    14,
				Forks:    1,
				Desc:     "Production SNMP daemon configuration for Raspberry Pi exposing hardware CPU core temperatures, system health, and Linux OS telemetry over network monitoring.",
				Tech:     []string{"SNMP", "Raspberry Pi", "Linux Sysadmin", "Bash"},
				URL:      "https://github.com/smford/rpi-snmpd-configuration",
			},
			{
				Name:     "Sunton-ESP32-2432S028R-Case",
				Category: "3DPrinting",
				Stars:    12,
				Forks:    0,
				Desc:     "Custom parametric 3D printable enclosure remix for the Sunton ESP32 2.8\" LCD + touch controller, optimized for M3 socket screws and heat-set inserts.",
				Tech:     []string{"3D Printing", "Bambu Lab", "PLA-CF", "CAD Remix", "ESP32-CYD"},
				URL:      "https://github.com/smford/Sunton-ESP32-2432S028R-Case",
			},
			{
				Name:     "eeh-esp32-rfid",
				Category: "Embedded",
				Stars:    11,
				Forks:    0,
				Desc:     "Industrial RFID access control and safety interlock system built for East Essex Hackspace to restrict unauthorized power to 3kW workshop tools.",
				Tech:     []string{"Embedded C++", "ESP32", "RC522 RFID", "Hardware Relays"},
				URL:      "https://github.com/smford/eeh-esp32-rfid",
			},
			{
				Name:     "bambu-poop-chute-magnet-remix",
				Category: "3DPrinting",
				Stars:    4,
				Forks:    0,
				Desc:     "Magnetic snap-lock purge chute remix for Bambu Lab 3D printers, featuring recessed neodymium magnet slots for quick filament purge removal.",
				Tech:     []string{"3D Printing", "CAD Remix", "Bambu Lab", "Magnets"},
				URL:      "https://github.com/smford/bambu-poop-chute-magnet-remix",
			},
			{
				Name:     "esp32-temp-monitor",
				Category: "Embedded",
				Stars:    2,
				Forks:    0,
				Desc:     "Autonomous IoT temperature regulation daemon with Dallas DS18B20 OneWire sensors, local SSR cooling relay control, and Pushover alerts.",
				Tech:     []string{"C++", "ESP32", "OneWire", "DS18B20", "Pushover API"},
				URL:      "https://github.com/smford/esp32-temp-monitor",
			},
			{
				Name:     "dhcp-routeros-tool",
				Category: "Networking",
				Stars:    1,
				Forks:    0,
				Desc:     "High-performance Go CLI utility querying MikroTik RouterOS v7 API to dump and inspect active DHCP leases with MAC address OUI vendor resolution.",
				Tech:     []string{"Go", "MikroTik API", "RouterOS v7", "Networking", "CLI"},
				URL:      "https://github.com/smford/dhcp-routeros-tool",
			},
			{
				Name:     "motioneye-snapshotter",
				Category: "RaspberryPi",
				Stars:    1,
				Forks:    0,
				Desc:     "Webhook-driven surveillance snapshot utility hooking into motionEye motion detection triggers to extract, timestamp, and archive security footage.",
				Tech:     []string{"Python", "motionEye", "Webhooks", "Raspberry Pi"},
				URL:      "https://github.com/smford/motioneye-snapshotter",
			},
			{
				Name:     "mdee",
				Category: "Terminal",
				Stars:    1,
				Forks:    0,
				Desc:     "High-fidelity terminal Markdown viewer in Go featuring full GFM rendering, Chroma syntax highlighting, OSC 1337 image support, and Mermaid diagram execution.",
				Tech:     []string{"Go", "Terminal CLI", "Goldmark", "Lip Gloss", "Chroma"},
				URL:      "https://github.com/smford/mdee",
				Demo:     "https://stephenford.org/mdee/",
			},
			{
				Name:     "delim",
				Category: "Terminal",
				Stars:    0,
				Forks:    0,
				Desc:     "Zero-allocation high-speed terminal delimiter and visual barrier generator in Go, engineered for SRE log-stream separation and build diagnostics.",
				Tech:     []string{"Go", "Terminal CLI", "Zero-Alloc", "Cross-Platform", "Homebrew"},
				URL:      "https://github.com/smford/delim",
			},
			{
				Name:     "hamburger-decomp",
				Category: "Web",
				Stars:    0,
				Forks:    0,
				Desc:     "Photorealistic, interactive 3D culinary exploded view featuring smooth scroll interpolation, anti-aliased composite layers, and procedural audio synthesis.",
				Tech:     []string{"JavaScript", "CSS 3D", "Web Audio API", "HTML5", "Vanilla Web"},
				URL:      "https://github.com/smford/hamburger-decomp",
				Demo:     "https://stephenford.org/hamburger-decomp/",
			},
			{
				Name:     "fly-around-the-moon",
				Category: "Web",
				Stars:    0,
				Forks:    0,
				Desc:     "Interactive 3D Moon exploration web application featuring dynamic Level-of-Detail (LOD), sub-kilometer surface photography, and real-time ephemeris.",
				Tech:     []string{"Three.js", "WebGL", "Vite", "JavaScript", "Astrodynamics"},
				URL:      "https://github.com/smford/fly-around-the-moon",
				Demo:     "https://stephenford.org/fly-around-the-moon/",
			},
			{
				Name:     "blackhole",
				Category: "Web",
				Stars:    0,
				Forks:    0,
				Desc:     "Chandra-X: Relativistic Schwarzschild black hole simulator featuring general relativistic raymarching, gravitational lensing, and Doppler beaming disk.",
				Tech:     []string{"WebGL 2.0", "GLSL Shaders", "General Relativity", "Raymarching"},
				URL:      "https://github.com/smford/blackhole",
				Demo:     "https://stephenford.org/blackhole/",
			},
			{
				Name:     "OpenLicenseGuide",
				Category: "Web",
				Stars:    0,
				Forks:    0,
				Desc:     "Interactive open-source license selector and decision tree designed with zero runtime dependencies to help developers choose the right OSS license in under 60 seconds.",
				Tech:     []string{"JavaScript", "SPDX", "PWA", "Offline-First", "Vanilla Web"},
				URL:      "https://github.com/smford/OpenLicenseGuide",
				Demo:     "https://stephenford.org/OpenLicenseGuide/",
			},
			{
				Name:     "vim-html",
				Category: "Terminal",
				Stars:    0,
				Forks:    0,
				Desc:     "Authentic in-browser Vim replica and simulated Linux SRE workstation environment featuring hierarchical VFS, GFM markdown rendering, and kernel boot simulation.",
				Tech:     []string{"JavaScript", "Vim Simulator", "Virtual Filesystem", "Terminal Shell"},
				URL:      "https://github.com/smford/vim-html",
				Demo:     "https://stephenford.org/vim-html/",
			},
			{
				Name:     "homebrew-tap",
				Category: "Terminal",
				Stars:    0,
				Forks:    0,
				Desc:     "Official Homebrew tap for smford CLI tools (cidr-calculator, mdee, delim, matrix-rain, matrix-cat) with automated multi-architecture formula packaging.",
				Tech:     []string{"Homebrew", "Formulas", "Ruby", "Bash", "macOS/Linux"},
				URL:      "https://github.com/smford/homebrew-tap",
				Demo:     "https://stephenford.org/homebrew-tap/",
			},
			{
				Name:     "matrix-rain",
				Category: "Terminal",
				Stars:    0,
				Forks:    0,
				Desc:     "High-performance authentic Matrix digital rain simulator for the terminal in Go, featuring half-width Katakana glyphs and 24-bit truecolor glowing heads.",
				Tech:     []string{"Go", "Terminal CLI", "VT100 / ANSI", "TrueColor", "Zero-Alloc"},
				URL:      "https://github.com/smford/matrix-rain",
				Demo:     "https://stephenford.org/matrix-rain/",
			},
			{
				Name:     "matrix-cat",
				Category: "Terminal",
				Stars:    0,
				Forks:    0,
				Desc:     "Digital rain terminal text and source code viewer (mcat) in Go, cascading characters down the screen with Chroma syntax highlighting.",
				Tech:     []string{"Go", "Terminal CLI", "Chroma Highlighting", "TrueColor", "Cascading Text"},
				URL:      "https://github.com/smford/matrix-cat",
				Demo:     "https://stephenford.org/matrix-cat/",
			},
			{
				Name:     "toolsplit",
				Category: "Web",
				Stars:    0,
				Forks:    0,
				Desc:     "Power tool promotional discount calculator and prorated return receipt engine optimizing Home Depot BOGO savings for Milwaukee, DeWalt, RIDGID, and Makita tooling.",
				Tech:     []string{"JavaScript", "CSS3", "Calculators", "Workshop Tooling", "Web App"},
				URL:      "https://github.com/smford/toolsplit",
				Demo:     "https://stephenford.org/toolsplit/",
			},
			{
				Name:     "rbac-editor",
				Category: "Web",
				Stars:    0,
				Forks:    0,
				Desc:     "100% client-side in-browser Kubernetes RBAC and YAML security validator, featuring real-time AST traversal, anchor/alias resolution, and dead-code analysis.",
				Tech:     []string{"TypeScript", "Kubernetes RBAC", "YAML AST", "Security ICE", "Client-Side"},
				URL:      "https://github.com/smford/rbac-editor",
				Demo:     "https://stephenford.org/rbac-editor/",
			},
			{
				Name:     "yaml-web-validator",
				Category: "Web",
				Stars:    0,
				Forks:    0,
				Desc:     "YAML Clean: Modern client-side YAML IDE and anchor manager for SREs to manage anchors (&), aliases (*), and merge keys (<<), scan for secrets, and diff dereferenced trees.",
				Tech:     []string{"TypeScript", "YAML AST", "PWA Offline", "SRE Tooling", "Zero-Telemetry"},
				URL:      "https://github.com/smford/yaml-web-validator",
				Demo:     "https://stephenford.org/yaml-web-validator/",
			},
			{
				Name:     "pcboard-web-emu",
				Category: "Web",
				Stars:    0,
				Forks:    0,
				Desc:     "High-fidelity retro BBS simulator recreating an IBM PC running MS-DOS 6.22 dialing into Clark Development PCBoard v15.22 with CP437 canvas rendering.",
				Tech:     []string{"TypeScript", "Canvas CP437", "Web Audio API", "MS-DOS BBS", "Retro Computing"},
				URL:      "https://github.com/smford/pcboard-web-emu",
				Demo:     "https://stephenford.org/pcboard-web-emu/",
			},
			{
				Name:     "golang-mermaid",
				Category: "Terminal",
				Stars:    0,
				Forks:    0,
				Desc:     "Production-grade Go module rendering high-resolution Mermaid diagrams in modern terminals via Kitty graphics protocol, iTerm2 OSC 1337, and DEC Sixel bitmap.",
				Tech:     []string{"Go", "Kitty Graphics", "iTerm2 OSC 1337", "DEC Sixel", "Mermaid CLI"},
				URL:      "https://github.com/smford/golang-mermaid",
			},
			{
				Name:     "cidr-calculator",
				Category: "Networking",
				Stars:    0,
				Forks:    0,
				Desc:     "High-performance zero-dependency IPv4/IPv6 CIDR range decomposition engine, route aggregator, collision detector, and capacity planner in Go.",
				Tech:     []string{"Go", "IPv4 & IPv6", "CIDR / Subnetting", "Route Aggregation", "Terraform", "SRE Tooling"},
				URL:      "https://github.com/smford/cidr-calculator",
				Demo:     "https://stephenford.org/cidr-calculator/",
			},
			{
				Name:     "emergency-stop-holder",
				Category: "3DPrinting",
				Stars:    0,
				Forks:    0,
				Desc:     "Industrial 3D printable wall mount enclosure for standard 22mm mushroom emergency stop switches, featuring tall (92mm) and short (75mm) models.",
				Tech:     []string{"3D Printing", "Workshop Safety", "CAD", "Bambu Lab", "Hardware"},
				URL:      "https://github.com/smford/emergency-stop-holder",
			},
			{
				Name:     "cladding-cutlist",
				Category: "Fabrication",
				Stars:    0,
				Forks:    0,
				Desc:     "Timber cladding cutlist generator and parametric layout optimizer minimizing board waste for workshop sheds and exterior timber architectural structures.",
				Tech:     []string{"JavaScript", "Timber Fabrication", "CAD", "Woodworking", "Web App"},
				URL:      "https://github.com/smford/cladding-cutlist",
			},
		},
	}
}
