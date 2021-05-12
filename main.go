module encrypt

import (
  "fmt"
  "bufio"
  "os"
  "math/rand"
  "time"
  "github.com/inancgumus/screen"
)

var (
  counter int
  keys []string
  line string
)

func keyGenerator(init int64) (keys []string) {
  rand.Seed(init)
  for x := 0; x < 2; x++ {
    bytes := make([]byte, 16)
    for i := 0; i < 16; i++ {
      bytes[i] = byte(65 + rand.Intn(25))
    }
    keys = append(keys, string(bytes))

    bytes = make([]byte, 16)
    for i := 0; i < 16; i++ {
      bytes[i] = byte(97 + rand.Intn(25))
    }
    keys = append(keys, string(bytes))

    bytes = make([]byte, 16)
    for i := 0; i < 16; i++ {
      bytes[i] = byte(48 + rand.Intn(9))
    }
    keys = append(keys, string(bytes))
  }
  return keys
}

func crypt_raw(input string, key string) (result []byte) {
  count := 0
  for i := 0; i < len(input); i++ {
    if count == len(key) {
      count = 0
    }
    result = append(result, byte(input[i]) ^ key[count])
    count++
	}
  
  return result
}

func crypt(input []byte, key string) (result []byte) {
  count := 0
  for i := 0; i < len(input); i++ {
    if count == len(key) {
		  count = 0
    }
    result = append(result, byte(input[i]) ^ key[count])
    count++
	}
  return result
}

func encrypt_file(name string) {
  file, encryption_error := os.Open(name)
  if encryption_error != nil {
    fmt.Println(encryption_error)
  }
  scanner := bufio.NewScanner(file)

  newFile, _ := os.OpenFile("encryptedFile.txt", os.O_CREATE|os.O_WRONLY, 700)
  newFile.Truncate(0)
  newFile.Seek(0,0)
  var result []byte
  rand.Seed(time.Now().UTC().UnixNano())
  init := int64(rand.Intn(255))
  keys := keyGenerator(init)
  result = append(result, byte(init))
  //fmt.Println(init)
  //fmt.Printf("%s\n%s\n%s\n", keys[0], keys[1], keys[2])
  _, err := newFile.Write(result)
  if err != nil {
    fmt.Println(err)
  }

  for scanner.Scan() {
    line = scanner.Text()
    
    first := crypt_raw(line, keys[0])
    second := crypt(first, keys[1])
    third := crypt(second, keys[2])
    fourth := crypt(third, keys[3])
    fifth := crypt(fourth, keys[4])
    result := crypt(fifth, keys[5])
    result = append(result, byte(10))

    _, err := newFile.Write(result)
    if err != nil {
      fmt.Println(err)
    }
  }
  file.Close()
  newFile.Close()
}

func encrypt_line(line string) (result []byte) {

  rand.Seed(time.Now().UTC().UnixNano())
  init := byte(rand.Intn(255))
  keys := keyGenerator(int64(init))
  first := crypt_raw(line, keys[0])
  second := crypt(first, keys[1])
  third := crypt(second, keys[2])
  fourth := crypt(third, keys[3])
  fifth := crypt(fourth, keys[4])
  result = crypt(fifth, keys[5])
  result = append(result, 0)
  copy(result[1:], result[0:])
  result[0] = init

  return result
}

func decrypt_file() {
  file, decryption_error := os.Open("encryptedFile.txt")
  if decryption_error != nil {
    fmt.Println(decryption_error)
  }
  scanner := bufio.NewScanner(file)

  newFile, _ := os.OpenFile("decryptedFile.txt", os.O_CREATE|os.O_WRONLY, 700)
  newFile.Truncate(0)
  newFile.Seek(0,0)
  var keys []string

  for scanner.Scan() {
    
    line = scanner.Text()
    
    if counter == 0 {
      init := byte(line[0])
      keys = keyGenerator(int64(init))
      first := crypt_raw(string(line[1:]), keys[5])
      second := crypt(first, keys[4])
      third := crypt(second, keys[3])
      fourth := crypt(third, keys[2])
      fifth := crypt(fourth, keys[1])
      result := crypt(fifth, keys[0])
      result = append(result, byte(10))
      _, err := newFile.Write(result)
      if err != nil {
        fmt.Println(err)
      }
      //fmt.Println(init)
      //fmt.Printf("%s\n%s\n%s\n", keys[0], keys[1], keys[2])
      counter = 1
    } else if counter != 0 {

      first := crypt_raw(line, keys[5])
      second := crypt(first, keys[4])
      third := crypt(second, keys[3])
      fourth := crypt(third, keys[2])
      fifth := crypt(fourth, keys[1])
      result := crypt(fifth, keys[0])
      result = append(result, byte(10))

      _, err := newFile.Write(result)
      if err != nil {
        fmt.Println(err)
      }
    }
  }
  file.Close()
  newFile.Close()
}

func decrypt_line(line []byte) (result []byte) {

  rand.Seed(time.Now().UTC().UnixNano())
  
  init := byte(line[0])
  keys = keyGenerator(int64(init))
  first := crypt_raw(string(line[1:]), keys[5])
  second := crypt(first, keys[4])
  third := crypt(second, keys[3])
  fourth := crypt(third, keys[2])
  fifth := crypt(fourth, keys[1])
  result = crypt(fifth, keys[0])

  return result
}

func main() {
  var input string
  screen.Clear()
  screen.MoveTopLeft()
  fmt.Printf("Welcome to EasyEncryption. Press:\n\n1) File encryption\n2) Line encryption\n\nInput: ")
  fmt.Scanln(&input)

  if input == "1" {
    fmt.Println("\nFile to encrypt:")
    fmt.Scanln(&input)
    encrypt_file(input)
    decrypt_file()
    fmt.Println("Done!")
  } else if input == "2" {
    fmt.Println("\nLine to encrypt:")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan() 
    input = scanner.Text()
    encrypted_result := encrypt_line(input)
    fmt.Printf("\nEncrypted: %s\n", encrypted_result)
    decrypted_result := decrypt_line(encrypted_result)
    fmt.Printf("\nDecrypted: %s\n", decrypted_result)
  }
}
