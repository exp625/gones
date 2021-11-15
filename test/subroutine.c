unsigned char x;
unsigned char val = 10;

void incX() {
    x++;
    if (x < val) {
        incX();
    }
    return;
}

int main() {
    x = 0;
    incX();
    // Return value will be stored in A register
    return x;
}

