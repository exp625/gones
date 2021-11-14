unsigned char main() {
    unsigned char i;
    unsigned char factorial = 1;
    for (i = 1; i <= 3; i++) {
    		factorial = factorial * i;
    }
    // Return value will be stored in A register
    return factorial;
}

