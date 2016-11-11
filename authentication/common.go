package authentication

//
//import (
//	"crypto/sha256"
//	"encoding/hex"
//	"os/user"
//)
//
//func getPassword(pass string, salt string) string {
//	sha_256 := sha256.New()
//	sha_256.Write([]byte(pass + salt))
//	return hex.EncodeToString(sha_256.Sum(nil))
//}
//
//func CheckPassword(usr user.User) (nil, errors error){
//	sha_256 := sha256.New()
//	sha_256.Write([]byte(lgn.Password + usr.Salt))
//	pass := hex.EncodeToString(sha_256.Sum(nil))
//
//	l_debug.Log("compare", usr.Password, pass)
//
//	// Check password
//	if usr.Password != pass {
//		return nil, errors.InvalidContent("Password does not match our records.")
//	}
//}
