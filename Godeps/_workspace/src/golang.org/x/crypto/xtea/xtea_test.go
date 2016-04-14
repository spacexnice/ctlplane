// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xtea

import (
    "testing"
)

// A sample test key for when we just want to initialize a cipher
var testKey = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

// Test that the block size for XTEA is correct
func TestBlocksize(t *testing.T) {
    if BlockSize != 8 {
        t.Errorf("BlockSize constant - expected 8, got %d", BlockSize)
        return
    }

    c, err := NewCipher(testKey)
    if err != nil {
        t.Errorf("NewCipher(%d bytes) = %s", len(testKey), err)
        return
    }

    result := c.BlockSize()
    if result != 8 {
        t.Errorf("BlockSize function - expected 8, got %d", result)
        return
    }
}

// A series of test values to confirm that the Cipher.table array was initialized correctly
var testTable = []uint32{
    0x00112233, 0x6B1568B8, 0xE28CE030, 0xC5089E2D, 0xC5089E2D, 0x1EFBD3A2, 0xA7845C2A, 0x78EF0917,
    0x78EF0917, 0x172682D0, 0x5B6AC714, 0x822AC955, 0x3DE68511, 0xDC1DFECA, 0x2062430E, 0x3611343F,
    0xF1CCEFFB, 0x900469B4, 0xD448ADF8, 0x2E3BE36D, 0xB6C46BF5, 0x994029F2, 0x994029F2, 0xF3335F67,
    0x6AAAD6DF, 0x4D2694DC, 0x4D2694DC, 0xEB5E0E95, 0x2FA252D9, 0x4551440A, 0x121E10D6, 0xB0558A8F,
    0xE388BDC3, 0x0A48C004, 0xC6047BC0, 0x643BF579, 0xA88039BD, 0x02736F32, 0x8AFBF7BA, 0x5C66A4A7,
    0x5C66A4A7, 0xC76AEB2C, 0x3EE262A4, 0x215E20A1, 0x215E20A1, 0x7B515616, 0x03D9DE9E, 0x1988CFCF,
    0xD5448B8B, 0x737C0544, 0xB7C04988, 0xDE804BC9, 0x9A3C0785, 0x3873813E, 0x7CB7C582, 0xD6AAFAF7,
    0x4E22726F, 0x309E306C, 0x309E306C, 0x8A9165E1, 0x1319EE69, 0xF595AC66, 0xF595AC66, 0x4F88E1DB,
}

// Test that the cipher context is initialized correctly
func TestCipherInit(t *testing.T) {
    c, err := NewCipher(testKey)
    if err != nil {
        t.Errorf("NewCipher(%d bytes) = %s", len(testKey), err)
        return
    }

    for i := 0; i < len(c.table); i++ {
        if c.table[i] != testTable[i] {
            t.Errorf("NewCipher() failed to initialize Cipher.table[%d] correctly. Expected %08X, got %08X", i, testTable[i], c.table[i])
            break
        }
    }
}

// Test that invalid key sizes return an error
func TestInvalidKeySize(t *testing.T) {
    // Test a long key
    key := []byte{
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
    }

    _, err := NewCipher(key)
    if err == nil {
        t.Errorf("Invalid key size %d didn't result in an error.", len(key))
    }

    // Test a short key
    key = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}

    _, err = NewCipher(key)
    if err == nil {
        t.Errorf("Invalid key size %d didn't result in an error.", len(key))
    }
}

// Test that we can correctly decode some bytes we have encoded
func TestEncodeDecode(t *testing.T) {
    original := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
    input := original
    output := make([]byte, BlockSize)

    c, err := NewCipher(testKey)
    if err != nil {
        t.Errorf("NewCipher(%d bytes) = %s", len(testKey), err)
        return
    }

    // Encrypt the input block
    c.Encrypt(output, input)

    // Check that the output does not match the input
    differs := false
    for i := 0; i < len(input); i++ {
        if output[i] != input[i] {
            differs = true
            break
        }
    }
    if differs == false {
        t.Error("Cipher.Encrypt: Failed to encrypt the input block.")
        return
    }

    // Decrypt the block we just encrypted
    input = output
    output = make([]byte, BlockSize)
    c.Decrypt(output, input)

    // Check that the output from decrypt matches our initial input
    for i := 0; i < len(input); i++ {
        if output[i] != original[i] {
            t.Errorf("Decrypted byte %d differed. Expected %02X, got %02X\n", i, original[i], output[i])
            return
        }
    }
}

// Test Vectors
type CryptTest struct {
    key        []byte
    plainText  []byte
    cipherText []byte
}

var CryptTests = []CryptTest{
    // These were sourced from http://www.freemedialibrary.com/index.php/XTEA_test_vectors
    {
        []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
        []byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48},
        []byte{0x49, 0x7d, 0xf3, 0xd0, 0x72, 0x61, 0x2c, 0xb5},
    },
    {
        []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
        []byte{0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41},
        []byte{0xe7, 0x8f, 0x2d, 0x13, 0x74, 0x43, 0x41, 0xd8},
    },
    {
        []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
        []byte{0x5a, 0x5b, 0x6e, 0x27, 0x89, 0x48, 0xd7, 0x7f},
        []byte{0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41},
    },
    {
        []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
        []byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48},
        []byte{0xa0, 0x39, 0x05, 0x89, 0xf8, 0xb8, 0xef, 0xa5},
    },
    {
        []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
        []byte{0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41},
        []byte{0xed, 0x23, 0x37, 0x5a, 0x82, 0x1a, 0x8c, 0x2d},
    },
    {
        []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
        []byte{0x70, 0xe1, 0x22, 0x5d, 0x6e, 0x4e, 0x76, 0x55},
        []byte{0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41},
    },

    // These vectors are from http://wiki.secondlife.com/wiki/XTEA_Strong_Encryption_Implementation#Bouncy_Castle_C.23_API
    {
        []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
        []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
        []byte{0xDE, 0xE9, 0xD4, 0xD8, 0xF7, 0x13, 0x1E, 0xD9},
    },
    {
        []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
        []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
        []byte{0x06, 0x5C, 0x1B, 0x89, 0x75, 0xC6, 0xA8, 0x16},
    },
    {
        []byte{0x01, 0x23, 0x45, 0x67, 0x12, 0x34, 0x56, 0x78, 0x23, 0x45, 0x67, 0x89, 0x34, 0x56, 0x78, 0x9A},
        []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
        []byte{0x1F, 0xF9, 0xA0, 0x26, 0x1A, 0xC6, 0x42, 0x64},
    },
    {
        []byte{0x01, 0x23, 0x45, 0x67, 0x12, 0x34, 0x56, 0x78, 0x23, 0x45, 0x67, 0x89, 0x34, 0x56, 0x78, 0x9A},
        []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
        []byte{0x8C, 0x67, 0x15, 0x5B, 0x2E, 0xF9, 0x1E, 0xAD},
    },
}

// Test encryption
func TestCipherEncrypt(t *testing.T) {
    for i, tt := range CryptTests {
        c, err := NewCipher(tt.key)
        if err != nil {
            t.Errorf("NewCipher(%d bytes), vector %d = %s", len(tt.key), i, err)
            continue
        }

        out := make([]byte, len(tt.plainText))
        c.Encrypt(out, tt.plainText)

        for j := 0; j < len(out); j++ {
            if out[j] != tt.cipherText[j] {
                t.Errorf("Cipher.Encrypt %d: out[%d] = %02X, expected %02X", i, j, out[j], tt.cipherText[j])
                break
            }
        }
    }
}

// Test decryption
func TestCipherDecrypt(t *testing.T) {
    for i, tt := range CryptTests {
        c, err := NewCipher(tt.key)
        if err != nil {
            t.Errorf("NewCipher(%d bytes), vector %d = %s", len(tt.key), i, err)
            continue
        }

        out := make([]byte, len(tt.cipherText))
        c.Decrypt(out, tt.cipherText)

        for j := 0; j < len(out); j++ {
            if out[j] != tt.plainText[j] {
                t.Errorf("Cipher.Decrypt %d: out[%d] = %02X, expected %02X", i, j, out[j], tt.plainText[j])
                break
            }
        }
    }
}
