#include <WiFi.h>
#include "DHT.h"
#include "time.h"
#include <LiquidCrystal.h>
#include <Wire.h>
#include <Adafruit_BMP085.h>
#include <HTTPClient.h>
#include <WiFiClientSecure.h>
#include <Ticker.h>

#define DHTPIN D11
#define DHTTYPE DHT11
#define GAS_DIN D12
#define GAS_AIN A0
#define BUZZER_PIN D13
#define LCD_BACKLIGHT A4
#define PUSH_BUTTON_PIN 1

int ad_value, gas;
float h, t, p;
const char *ssid = "esp32-wifi";
const char *password = "123456789";
const char *serverName = "https://duthweather.azurewebsites.net/api/add";
const char *ntpServer = "pool.ntp.org";
const long gmtOffset_sec = 7200;
const int daylightOffset_sec = 3600;
unsigned long lastDataSentTime = 0;
const unsigned long sendDataInterval = 2 * 60 * 1000;
volatile bool backlightState = false;

Ticker gasSensorTicker;
Ticker buttonCheckTicker;
LiquidCrystal lcd(26, 2, 17, 14, 13, 25);  // ( RS, EN, D4, D5, D6, D7 )
DHT dht(DHTPIN, DHTTYPE);
Adafruit_BMP085 bmp;

void IRAM_ATTR handleButtonPress() {
  backlightState = !backlightState;
  digitalWrite(LCD_BACKLIGHT, backlightState ? HIGH : LOW);
}

void gasSensorCheck() {
  ad_value = analogRead(GAS_AIN);
  if (ad_value >= 2500) {
    while (ad_value >= 2500) {
      lcd.clear();
      delay(50);
      lcd.setCursor(0, 0);
      lcd.print("DANGER! Gas");
      lcd.setCursor(0, 1);
      lcd.print("leakage detected");
      digitalWrite(BUZZER_PIN, HIGH);
      delay(100);
      digitalWrite(BUZZER_PIN, LOW);
      delay(100);
      ad_value = analogRead(GAS_AIN);
      delay(50);
    }
  }
}

void checkButton() {
  static bool lastButtonState = HIGH;
  bool currentButtonState = digitalRead(PUSH_BUTTON_PIN);
  if (lastButtonState == HIGH && currentButtonState == LOW) {
    handleButtonPress();
  }
  lastButtonState = currentButtonState;
}

const char *rootCACertificate =
  "-----BEGIN CERTIFICATE-----\n"
  "MIIDjjCCAnagAwIBAgIQAzrx5qcRqaC7KGSxHQn65TANBgkqhkiG9w0BAQsFADBh\n"
  "MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\n"
  "d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH\n"
  "MjAeFw0xMzA4MDExMjAwMDBaFw0zODAxMTUxMjAwMDBaMGExCzAJBgNVBAYTAlVT\n"
  "MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdpY2VydC5j\n"
  "b20xIDAeBgNVBAMTF0RpZ2lDZXJ0IEdsb2JhbCBSb290IEcyMIIBIjANBgkqhkiG\n"
  "9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuzfNNNx7a8myaJCtSnX/RrohCgiN9RlUyfuI\n"
  "2/Ou8jqJkTx65qsGGmvPrC3oXgkkRLpimn7Wo6h+4FR1IAWsULecYxpsMNzaHxmx\n"
  "1x7e/dfgy5SDN67sH0NO3Xss0r0upS/kqbitOtSZpLYl6ZtrAGCSYP9PIUkY92eQ\n"
  "q2EGnI/yuum06ZIya7XzV+hdG82MHauVBJVJ8zUtluNJbd134/tJS7SsVQepj5Wz\n"
  "tCO7TG1F8PapspUwtP1MVYwnSlcUfIKdzXOS0xZKBgyMUNGPHgm+F6HmIcr9g+UQ\n"
  "vIOlCsRnKPZzFBQ9RnbDhxSJITRNrw9FDKZJobq7nMWxM4MphQIDAQABo0IwQDAP\n"
  "BgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQEAwIBhjAdBgNVHQ4EFgQUTiJUIBiV\n"
  "5uNu5g/6+rkS7QYXjzkwDQYJKoZIhvcNAQELBQADggEBAGBnKJRvDkhj6zHd6mcY\n"
  "1Yl9PMWLSn/pvtsrF9+wX3N3KjITOYFnQoQj8kVnNeyIv/iPsGEMNKSuIEyExtv4\n"
  "NeF22d+mQrvHRAiGfzZ0JFrabA0UWTW98kndth/Jsw1HKj2ZL7tcu7XUIOGZX1NG\n"
  "Fdtom/DzMNU+MeKNhJ7jitralj41E6Vf8PlwUHBHQRFXGU7Aj64GxJUTFy8bJZ91\n"
  "8rGOmaFvE7FBcf6IKshPECBV1/MUReXgRPTqh5Uykw7+U0b6LJ3/iyK5S9kJRaTe\n"
  "pLiaWN0bfVKfjllDiIGknibVb63dDcY3fe0Dkhvld1927jyNxF1WW6LZZm6zNTfl\n"
  "MrY=\n"
  "-----END CERTIFICATE-----\n";

String printLocalTime() {
  struct tm timeinfo;
  if (!getLocalTime(&timeinfo)) {
    return "";
  }
  time_t timestamp = mktime(&timeinfo);
  char formattedTime[30];
  strftime(formattedTime, sizeof(formattedTime), "%Y-%m-%dT%H:%M:%SZ", &timeinfo);
  return String(formattedTime);
}

void setup() {
  lcd.begin(16, 2);
  dht.begin();
  if (!bmp.begin()) {
    while (1) {}
  }

  pinMode(GAS_DIN, INPUT);
  pinMode(GAS_AIN, INPUT);
  pinMode(PUSH_BUTTON_PIN, INPUT);
  pinMode(BUZZER_PIN, OUTPUT);
  pinMode(LCD_BACKLIGHT, OUTPUT);
  digitalWrite(LCD_BACKLIGHT, HIGH);

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Connecting to");
  lcd.setCursor(0, 1);
  lcd.print("wifi");
  WiFi.begin(ssid, password);

  unsigned long startAttemptTime = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - startAttemptTime < 10000) {
    for (int i = 0; i < 3; i++) {
      delay(500);
      lcd.print(".");
    }
    delay(500);
    lcd.setCursor(0, 1);
    lcd.print("wifi   ");
    delay(500);
    lcd.setCursor(0, 1);
    lcd.print("wifi");
  }

  if (WiFi.status() == WL_CONNECTED) {
    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print("Connected to:");
    lcd.setCursor(0, 1);
    lcd.print(ssid);

    configTime(gmtOffset_sec, daylightOffset_sec, ntpServer);
    printLocalTime();
  } else {
    lcd.clear();
    lcd.home();
    lcd.print("Couldn't connect");
    lcd.setCursor(0, 1);
    lcd.print("to wifi");
    delay(2500);
    lcd.clear();
  }

  buttonCheckTicker.attach(0.05, checkButton);
  gasSensorTicker.attach(0.5, gasSensorCheck);
}

void DisplaySensorData(int ad_value, float h, float t, float p) {
  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Humidity: ");
  lcd.setCursor(0, 1);
  lcd.print(h);
  lcd.print(" %");
  delay(2500);

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Temperature: ");
  lcd.setCursor(0, 1);
  lcd.print(t);
  lcd.print(" C");
  delay(2500);

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Gas sensor: ");
  lcd.setCursor(0, 1);
  lcd.print(ad_value);
  lcd.print(" ppm");
  delay(2500);

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Pressure: ");
  lcd.setCursor(0, 1);
  lcd.print(p);
  lcd.print(" hPa");
  delay(2500);
}

String createPayload(const String &sensorName, const String &timestamp, const String &data) {
  return "{\"SensorName\":\"" + sensorName + "\",\"SensorData\":{\"timestamp\":\"" + timestamp + "\"," + data + "}}";
}

void sendData(int ad_value, float h, float t, float p) {
  const int maxRetries = 3;
  int attempts = 0;
  bool success = false;

  while (attempts < maxRetries && !success) {
    WiFiClientSecure *client = new WiFiClientSecure;
    if (client) {
      lcd.clear();
      lcd.setCursor(0, 0);
      lcd.print("Uploading data");
      delay(500);

      client->setCACert(rootCACertificate);

      HTTPClient https;
      if (https.begin(*client, serverName)) {
        String timestamp = printLocalTime();
        int httpResponseCode;

        String payload_dht11 = createPayload("dht11", timestamp, "\"humidity\":" + String(h));
        httpResponseCode = https.POST(payload_dht11);

        if (httpResponseCode > 0) {
          String payload_bmp180 = createPayload("bmp180", timestamp, "\"temperature\":" + String(t) + ",\"pressure\":" + String(p));
          httpResponseCode = https.POST(payload_bmp180);

          if (httpResponseCode > 0) {
            String payload_mq135 = createPayload("mq135", timestamp, "\"gas_level\":" + String(ad_value));
            httpResponseCode = https.POST(payload_mq135);
          }
        }

        https.end();

        if (httpResponseCode > 0) {
          success = true;
          lcd.setCursor(0, 1);
          lcd.print("Done!");
          delay(2000);
          lcd.clear();
        } else {
          attempts++;
          lcd.setCursor(0, 1);
          lcd.print("Retrying...");
          delay(2500);
          lcd.clear();
        }
      } else {
        lcd.setCursor(0, 1);
        lcd.print("Connection failed");
        delay(2500);
        lcd.clear();
        attempts++;
      }
    }
    delete client;
  }

  if (!success) {
    lcd.setCursor(0, 1);
    lcd.print("Failed after retries");
    delay(2000);
    lcd.clear();
  }
}

void loop() {
  unsigned long currentTime = millis();

  if (currentTime - lastDataSentTime >= sendDataInterval) {
    lastDataSentTime = currentTime;
    // ad_value = analogRead(GAS_AIN);
    gas = analogRead(GAS_AIN);
    ad_value = gas - 120;
    ad_value = map(ad_value, 0, 1024, 400, 5000);
    h = dht.readHumidity();
    t = bmp.readTemperature();
    p = bmp.readPressure() / 100.0;
    if (WiFi.status() == WL_CONNECTED) {
      sendData(ad_value, h, t, p);
    }
  } else {
    // ad_value = analogRead(GAS_AIN);
    gas = analogRead(GAS_AIN);
    ad_value = gas - 120;
    ad_value = map(ad_value, 0, 1024, 400, 5000);
    h = dht.readHumidity();
    t = bmp.readTemperature();
    p = bmp.readPressure() / 100.0;
    DisplaySensorData(ad_value, h, t, p);
  }
}