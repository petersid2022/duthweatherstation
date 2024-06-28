#include <WiFi.h>
#include "DHT.h"
#include "time.h"
#include <LiquidCrystal.h>
#include <Wire.h>
#include <Adafruit_BMP085.h>
#include <WiFiClientSecure.h>
#include <Ticker.h>

#define DHTPIN 11
#define DHTTYPE DHT11
#define GAS_DIN 12
#define GAS_AIN A0
#define BUZZER_PIN 13
#define LCD_BACKLIGHT A4
#define PUSH_BUTTON_PIN 1

int ad_value, gas;
float h, t, p;
const char *ssid = "esp32-wifi";
const char *password = "123456789";
const char *ntpServer = "pool.ntp.org";
const long gmtOffset_sec = 7200;
const int daylightOffset_sec = 3600;
volatile bool backlightState = false;
bool alarmTriggeredToday = false;

Ticker buttonCheckTicker;
LiquidCrystal lcd(26, 2, 17, 14, 13, 25);
DHT dht(DHTPIN, DHTTYPE);
Adafruit_BMP085 bmp;

void IRAM_ATTR handleButtonPress() {
  backlightState = !backlightState;
  digitalWrite(LCD_BACKLIGHT, backlightState ? HIGH : LOW);
}

void checkButton() {
  static bool lastButtonState = HIGH;
  bool currentButtonState = digitalRead(PUSH_BUTTON_PIN);
  if (lastButtonState == HIGH && currentButtonState == LOW) {
    handleButtonPress();
  }
  lastButtonState = currentButtonState;
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
  pinMode(8, OUTPUT);
  digitalWrite(8, LOW);

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
  } else {
    lcd.clear();
    lcd.home();
    lcd.print("Couldn't connect");
    lcd.setCursor(0, 1);
    lcd.print("to wifi");
  }
  delay(2500);
  buttonCheckTicker.attach(0.05, checkButton);
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

void AlarmClock() {
  struct tm timeinfo;
  if (!getLocalTime(&timeinfo)) {
    return;
  }

  int currentHour = timeinfo.tm_hour;
  int currentMinute = timeinfo.tm_min;

  if ((currentHour == 15 && currentMinute == 5) && !alarmTriggeredToday) {  // Set alarm time here
    alarmTriggeredToday = true;

    digitalWrite(LCD_BACKLIGHT, HIGH);
    lcd.clear();

    unsigned long startMillis = millis();
    unsigned long currentMillis = startMillis;

    while (currentMillis - startMillis < 30000) {  // Alarm duration (30 seconds)
      lcd.setCursor(0, 0);
      lcd.print("Good morning!");
      lcd.setCursor(0, 1);
      lcd.print("Time to wake up!");
      digitalWrite(BUZZER_PIN, HIGH);
      delay(150);
      digitalWrite(BUZZER_PIN, LOW);
      delay(150);
      if (touchRead(T1) <= 40) {  // Assuming a touch sensor on pin T1
        return;
      }
      currentMillis = millis();
    }
  }

  if (currentHour == 0 && currentMinute == 0) {
    alarmTriggeredToday = false;
  }
}

void loop() {
  AlarmClock();
  gas = analogRead(GAS_AIN);
  ad_value = gas - 120;
  ad_value = map(ad_value, 0, 1024, 400, 5000);
  h = dht.readHumidity();
  t = bmp.readTemperature();
  p = bmp.readPressure() / 100.0;
  if (isnan(h) || isnan(t) || isnan(ad_value)) {
    return;
  }
  DisplaySensorData(ad_value, h, t, p);
}
