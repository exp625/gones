unsigned char main() {
    unsigned char i;
    unsigned char factorial = 1;
    for (i = 1; i <= 3; i++) {
    		factorial = factorial * i;
    }
    // Return value will be stored in A register
    *(unsigned char*) 0x00F0 = factorial;
    return factorial;
}

