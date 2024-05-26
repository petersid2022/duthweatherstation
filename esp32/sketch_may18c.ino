#include "DHT.h"
#include "time.h"
#include <LiquidCrystal.h>
#include <Wire.h>
#include <Adafruit_BMP085.h>
#include <HTTPClient.h>
#include <WiFiClientSecure.h>

#define DHTPIN D11
#define DHTTYPE DHT11
#define GAS_DIN D12
#define GAS_AIN A0
#define BUZZER_PIN D13
#define LCD_BACKLIGHT A4

// ( RS, EN, D4, D5, D6, D7 )
LiquidCrystal lcd(26, 2, 17, 14, 13, 25);
DHT dht(DHTPIN, DHTTYPE);
Adafruit_BMP085 bmp;

const char *ssid = "tzampa";
const char *password = "geiasoupetro";
const char *serverName = "https://duthweatherstation.fly.dev/api/add";
const char* ntpServer = "pool.ntp.org";
const long  gmtOffset_sec = 7200;
const int   daylightOffset_sec = 3600;

const char* rootCACertificate = \
                                "-----BEGIN CERTIFICATE-----\n" \
                                "MIIFazCCA1OgAwIBAgIRAIIQz7DSQONZRGPgu2OCiwAwDQYJKoZIhvcNAQELBQAw\n" \
                                "TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh\n" \
                                "cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMTUwNjA0MTEwNDM4\n" \
                                "WhcNMzUwNjA0MTEwNDM4WjBPMQswCQYDVQQGEwJVUzEpMCcGA1UEChMgSW50ZXJu\n" \
                                "ZXQgU2VjdXJpdHkgUmVzZWFyY2ggR3JvdXAxFTATBgNVBAMTDElTUkcgUm9vdCBY\n" \
                                "MTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAK3oJHP0FDfzm54rVygc\n" \
                                "h77ct984kIxuPOZXoHj3dcKi/vVqbvYATyjb3miGbESTtrFj/RQSa78f0uoxmyF+\n" \
                                "0TM8ukj13Xnfs7j/EvEhmkvBioZxaUpmZmyPfjxwv60pIgbz5MDmgK7iS4+3mX6U\n" \
                                "A5/TR5d8mUgjU+g4rk8Kb4Mu0UlXjIB0ttov0DiNewNwIRt18jA8+o+u3dpjq+sW\n" \
                                "T8KOEUt+zwvo/7V3LvSye0rgTBIlDHCNAymg4VMk7BPZ7hm/ELNKjD+Jo2FR3qyH\n" \
                                "B5T0Y3HsLuJvW5iB4YlcNHlsdu87kGJ55tukmi8mxdAQ4Q7e2RCOFvu396j3x+UC\n" \
                                "B5iPNgiV5+I3lg02dZ77DnKxHZu8A/lJBdiB3QW0KtZB6awBdpUKD9jf1b0SHzUv\n" \
                                "KBds0pjBqAlkd25HN7rOrFleaJ1/ctaJxQZBKT5ZPt0m9STJEadao0xAH0ahmbWn\n" \
                                "OlFuhjuefXKnEgV4We0+UXgVCwOPjdAvBbI+e0ocS3MFEvzG6uBQE3xDk3SzynTn\n" \
                                "jh8BCNAw1FtxNrQHusEwMFxIt4I7mKZ9YIqioymCzLq9gwQbooMDQaHWBfEbwrbw\n" \
                                "qHyGO0aoSCqI3Haadr8faqU9GY/rOPNk3sgrDQoo//fb4hVC1CLQJ13hef4Y53CI\n" \
                                "rU7m2Ys6xt0nUW7/vGT1M0NPAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNV\n" \
                                "HRMBAf8EBTADAQH/MB0GA1UdDgQWBBR5tFnme7bl5AFzgAiIyBpY9umbbjANBgkq\n" \
                                "hkiG9w0BAQsFAAOCAgEAVR9YqbyyqFDQDLHYGmkgJykIrGF1XIpu+ILlaS/V9lZL\n" \
                                "ubhzEFnTIZd+50xx+7LSYK05qAvqFyFWhfFQDlnrzuBZ6brJFe+GnY+EgPbk6ZGQ\n" \
                                "3BebYhtF8GaV0nxvwuo77x/Py9auJ/GpsMiu/X1+mvoiBOv/2X/qkSsisRcOj/KK\n" \
                                "NFtY2PwByVS5uCbMiogziUwthDyC3+6WVwW6LLv3xLfHTjuCvjHIInNzktHCgKQ5\n" \
                                "ORAzI4JMPJ+GslWYHb4phowim57iaztXOoJwTdwJx4nLCgdNbOhdjsnvzqvHu7Ur\n" \
                                "TkXWStAmzOVyyghqpZXjFaH3pO3JLF+l+/+sKAIuvtd7u+Nxe5AW0wdeRlN8NwdC\n" \
                                "jNPElpzVmbUq4JUagEiuTDkHzsxHpFKVK7q4+63SM1N95R1NbdWhscdCb+ZAJzVc\n" \
                                "oyi3B43njTOQ5yOf+1CceWxG1bQVs5ZufpsMljq4Ui0/1lvh+wjChP4kqKOJ2qxq\n" \
                                "4RgqsahDYVvTH9w7jXbyLeiNdd8XM2w9U/t7y0Ff/9yi0GE44Za4rF2LN9d11TPA\n" \
                                "mRGunUHBcnWEvgJBQl9nJEiU0Zsnvgc/ubhPgXRR4Xq37Z0j4r7g1SgEEzwxA57d\n" \
                                "emyPxgcYxn/eR44/KJ4EBs+lVDR3veyJm+kXQ99b21/+jh5Xos1AnX5iItreGCc=\n" \
                                "-----END CERTIFICATE-----\n";

String printLocalTime() {
  struct tm timeinfo;
  if (!getLocalTime(&timeinfo)) {
    // Serial.println("Failed to obtain time");
    return "";
  }

  // Convert struct tm to time_t
  time_t timestamp = mktime(&timeinfo);

  // Format the timestamp according to RFC 3339
  char formattedTime[30];
  strftime(formattedTime, sizeof(formattedTime), "%Y-%m-%dT%H:%M:%SZ", &timeinfo);

  return String(formattedTime);
}

void setup() {
  // Serial.begin(115200);
  lcd.begin(16, 2);
  lcd.clear();
  digitalWrite(LCD_BACKLIGHT, HIGH);
  dht.begin();

  if (!bmp.begin()) {
    while (1) {}
  }

  pinMode(GAS_DIN, INPUT);
  pinMode(GAS_AIN, INPUT);
  pinMode(BUZZER_PIN, OUTPUT);
  pinMode(LCD_BACKLIGHT, OUTPUT);

  // Connect to WiFi
  // Serial.print("Attempting to connect to SSID: ");
  // Serial.println(ssid);
  WiFi.begin(ssid, password);

  // attempt to connect to Wifi network:
  while (WiFi.status() != WL_CONNECTED) {
    // Serial.print(".");
    delay(1000);
  }

  // Serial.print("Connected to ");
  // Serial.println(ssid);
  // Serial.println(WiFi.localIP());
  lcd.setCursor(0, 0);
  lcd.print("Connected to");
  lcd.setCursor(0, 1);
  lcd.print(ssid);

  // Init and get the time
  configTime(gmtOffset_sec, daylightOffset_sec, ntpServer);
  printLocalTime();
}

void loop() {
  lcd.clear();
  // Read values from sensors
  int ad_value = analogRead(GAS_AIN);
  float h = dht.readHumidity();
  float t = dht.readTemperature();
  float p = bmp.readPressure() / 100.0;

  if (ad_value >= 1800) {
    while (ad_value >= 1800) {
      lcd.clear();
      lcd.setCursor(0, 0);
      lcd.print("DANGER! Gas");
      lcd.setCursor(0, 1);
      lcd.print("leakage detected");
      digitalWrite(BUZZER_PIN, HIGH);
      delay(100);
      digitalWrite(BUZZER_PIN, LOW);
      delay(100);

      // Update ad_value within the loop
      ad_value = analogRead(GAS_AIN);
    }
  }

  // Display sensor readings
  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Humidity: ");
  lcd.setCursor(0, 1);
  lcd.print(h);
  lcd.print(" %");
  delay(2000);

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Temperature: ");
  lcd.setCursor(0, 1);
  lcd.print(t);
  lcd.print(" C");
  delay(2000);

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Gas sensor: ");
  lcd.setCursor(0, 1);
  lcd.print(ad_value);
  delay(2000);

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Pressure: ");
  lcd.setCursor(0, 1);
  lcd.print(p);
  lcd.print(" hPa");
  delay(2000);

  WiFiClientSecure *client = new WiFiClientSecure;
  if (client) {
    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print("Sending data...");
    
    client->setCACert(rootCACertificate);

    HTTPClient https;

    // Serial.print("[HTTPS] begin...\n");
    if (https.begin(*client, serverName)) {
      // Serial.print("[HTTPS] POST...\n");

      // Create JSON payload for DHT11
      String payload_dht11 = "{";
      payload_dht11 += "\"SensorName\":\"dht11\",";
      payload_dht11 += "\"SensorData\":{";
      payload_dht11 += "\"timestamp\":\"" + printLocalTime() + "\",";
      payload_dht11 += "\"temperature\":" + String(t) + ",";
      payload_dht11 += "\"humidity\":" + String(h);
      payload_dht11 += "}}";
      // Serial.println("Sending DHT11 data: " + payload_dht11);

      // start connection and send HTTP header
      int httpCode = https.POST(payload_dht11);

      // httpCode will be negative on error
      if (httpCode > 0) {
        // Serial.printf("[HTTPS] POST... code: %d\n", httpCode);
        // file found at server
        if (httpCode == HTTP_CODE_OK || httpCode == HTTP_CODE_CREATED) {
          // print server response payload
          String payload_dht11 = https.getString();
          // Serial.println(payload_dht11);
        }
      }

      // Create JSON payload for BMP180
      String payload_bmp180 = "{";
      payload_bmp180 += "\"SensorName\":\"bmp180\",";
      payload_bmp180 += "\"SensorData\":{";
      payload_bmp180 += "\"timestamp\":\"" + printLocalTime() + "\",";
      payload_bmp180 += "\"pressure\":" + String(p);
      payload_bmp180 += "}}";
      // Serial.println("Sending BMP180 data: " + payload_bmp180);

      // start connection and send HTTP header
      httpCode = https.POST(payload_bmp180);

      // httpCode will be negative on error
      if (httpCode > 0) {
        // Serial.printf("[HTTPS] POST... code: %d\n", httpCode);
        // file found at server
        if (httpCode == HTTP_CODE_OK || httpCode == HTTP_CODE_CREATED) {
          // print server response payload
          String payload_bmp180 = https.getString();
          // Serial.println(payload_bmp180);
        }
      }

      // Create JSON payload for MQ135
      String payload_mq135 = "{";
      payload_mq135 += "\"SensorName\":\"mq135\",";
      payload_mq135 += "\"SensorData\":{";
      payload_mq135 += "\"timestamp\":\"" + printLocalTime() + "\",";
      payload_mq135 += "\"gas_level\":" + String(ad_value);
      payload_mq135 += "}}";
      // Serial.println("Sending MQ135 data: " + payload_mq135 );

      // start connection and send HTTP header
      httpCode = https.POST(payload_mq135);

      // httpCode will be negative on error
      if (httpCode > 0) {
        // Serial.printf("[HTTPS] POST... code: %d\n", httpCode);
        // file found at server
        if (httpCode == HTTP_CODE_OK || httpCode == HTTP_CODE_CREATED) {
          // print server response payload
          String payload_mq135 = https.getString();
          // Serial.println(payload_mq135);
        }
      }

      https.end();
      lcd.setCursor(0, 1);
      lcd.print("Done!");
      delay(1000);
      lcd.clear();
    }
    // else {
    //   Serial.printf("[HTTPS] Unable to connect\n");
    // }
  }
  // else {
  //   Serial.printf("[HTTPS] Unable to create client\n");
  // }

  // digitalWrite(LCD_BACKLIGHT, LOW);
}
