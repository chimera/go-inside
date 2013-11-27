const int LOCK = 7;
const int DELAY = 1000;
const int BAUD = 19200;

void setup() {
  Serial.begin(BAUD);
  pinMode(LOCK, OUTPUT);
}

void loop() {
  while (Serial.available() > 0) {
    int val = Serial.parseInt();
    Serial.write("val received");
    if (val == 1) {
      Serial.write("unlock");
      unlock();
    }
  }
}

void unlock() {
  digitalWrite(LOCK, HIGH);
  delay(DELAY);
  digitalWrite(LOCK, LOW);
}
