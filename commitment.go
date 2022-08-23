//Using Schnorr to commit and decommit
//Make changes to the commit and decommit function according to your need i.e. take file names as arguments for the function

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/dedis/kyber.v2"
	//"gopkg.in/dedis/kyber.v2/group/edwards25519"
	"gopkg.in/dedis/kyber.v2/util/encoding"
)

type Data struct {
	s       string
	pub     string
	message string
}

// type Signature struct {
// 	r kyber.Point
// 	s kyber.Scalar
// }

// func Hash(s string) kyber.Scalar {
// 	sha256.Reset()
// 	sha256.Write([]byte(s))

// 	return curve.Scalar().SetBytes(sha256.Sum(nil))
// }

// m: Message
// x: Private key
// func Sign(m string, x kyber.Scalar) Signature {
// 	// Get the base of the curve.
// 	g := curve.Point().Base()

// 	// Pick a random k from allowed set.
// 	k := curve.Scalar().Pick(curve.RandomStream())

// 	// r = k * G (a.k.a the same operation as r = g^k)
// 	r := curve.Point().Mul(k, g)

// 	// Hash(m || r)
// 	e := Hash(m + r.String())

// 	// s = k - e * x
// 	s := curve.Scalar().Sub(k, curve.Scalar().Mul(e, x))

// 	return Signature{r: r, s: s}
// }

// m: Message
// S: Signature
func Comit_PublicKey(m string, S Signature) kyber.Point {
	// Create a generator.
	g := curve.Point().Base()

	// e = Hash(m || r)
	e := Hash(m + S.r.String())

	// y = (r - s * G) * (1 / e)
	y := curve.Point().Sub(S.r, curve.Point().Mul(S.s, g))
	y = curve.Point().Mul(curve.Scalar().Div(curve.Scalar().One(), e), y)

	return y
}

// m: Message
// s: Signature
// y: Public key
func Comit_Verify(m string, S Signature, y kyber.Point) bool {
	// Create a generator.
	g := curve.Point().Base()

	// e = Hash(m || r)
	e := Hash(m + S.r.String())

	// Attempt to reconstruct 's * G' with a provided signature; s * G = r - e * y
	sGv := curve.Point().Sub(S.r, curve.Point().Mul(e, y))

	// Construct the actual 's * G'
	sG := curve.Point().Mul(S.s, g)

	//fmt.Println(sG)
	//fmt.Println(sGv)
	// Equality check; ensure signature and public key outputs to s * G.
	return sG.Equal(sGv)
}

func (S Signature) String() string {
	return fmt.Sprintf("(r=%s, s=%s)", S.r, S.s)
}

func Commitment(x kyber.Scalar, m string, peer_number string) {
	path1 := "Commitment/" + peer_number + "/KGC"
	err := os.MkdirAll(path1, os.ModePerm)
	if err != nil {
		panic(err)
	}
	publicKey := curve.Point().Mul(x, curve.Point().Base())
	sig := Sign(m, x)

	f1, e1 := os.OpenFile(path1+"/Signature_S.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e1 != nil {
		fmt.Println(e1)
	}
	encoding.WriteHexScalar(curve, f1, sig.s)

	f2, e2 := os.OpenFile(path1+"/PubKey.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e2 != nil {
		fmt.Println(e2)
	}
	encoding.WriteHexPoint(curve, f2, publicKey)

	f3, e3 := os.OpenFile(path1+"/Message.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e3 != nil {
		fmt.Println(e3)
	}
	f3.WriteString(m)
	f3.Close()
	f4, e4 := os.OpenFile("Commitment/"+peer_number+"/KGD.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e4 != nil {
		fmt.Println(e4)
	}
	encoding.WriteHexPoint(curve, f4, sig.r)
	fmt.Printf("Commitment Done for Peer %s \n", peer_number)
}
func Commitment_sign(x kyber.Scalar, m string, peer_number string) {
	path1 := "Commitment/Signing/" + peer_number + "/KGC"
	err := os.MkdirAll(path1, os.ModePerm)
	if err != nil {
		panic(err)
	}
	publicKey := curve.Point().Mul(x, curve.Point().Base())
	sig := Sign(m, x)

	f1, e1 := os.OpenFile(path1+"/Signature_S.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e1 != nil {
		fmt.Println(e1)
	}
	encoding.WriteHexScalar(curve, f1, sig.s)

	f2, e2 := os.OpenFile(path1+"/PubKey.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e2 != nil {
		fmt.Println(e2)
	}
	encoding.WriteHexPoint(curve, f2, publicKey)

	f3, e3 := os.OpenFile(path1+"/Message.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e3 != nil {
		fmt.Println(e3)
	}
	f3.WriteString(m)
	f3.Close()
	f4, e4 := os.OpenFile("Commitment/Signing/"+peer_number+"/KGD.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if e4 != nil {
		fmt.Println(e4)
	}
	encoding.WriteHexPoint(curve, f4, sig.r)
	fmt.Printf("Sign Commitment Done for Peer %s \n", peer_number)
}

func Decommitment_j(peer_number string) string {
	path := "Broadcast/" + peer_number
	f1, e1 := os.Open(path + "/Signature_S.txt")
	if e1 != nil {
		fmt.Println(e1)
	}
	sig_d, e := encoding.ReadHexScalar(curve, f1)
	if e != nil {
		fmt.Println(e)
	}

	f2, e2 := os.Open(path + "/PubKey.txt")
	if e2 != nil {
		fmt.Println(e2)
	}
	pub_key, e_2 := encoding.ReadHexPoint(curve, f2)
	if e_2 != nil {
		fmt.Println(e_2)
	}
	path1 := "Broadcast/" + peer_number + "/KGD.txt"
	f3, e3 := os.Open(path1)
	if e3 != nil {
		fmt.Println(e3)
	}
	KGD_j, e_3 := encoding.ReadHexPoint(curve, f3)
	if e_3 != nil {
		fmt.Println(e_3)
	}

	message, e4 := ioutil.ReadFile(path + "/Message.txt")
	if e4 != nil {
		fmt.Println(e4)
	}

	newS := Signature{}
	newS.s = sig_d
	newS.r = KGD_j
	// fmt.Println(string(message))
	// fmt.Println(pub_key)
	// fmt.Println(sig_d)
	//fmt.Println(newS.s)
	t := Comit_Verify(string(message), newS, pub_key)
	//fmt.Println(t)
	if t {

		return pub_key.String()
	} else {
		return "Invalid"
	}
}

func Decommitment_j_sign(peer_number string) string {
	path := "Broadcast/" + peer_number + "/Signing"
	f1, e1 := os.Open(path + "/Signature_S.txt")
	if e1 != nil {
		fmt.Println(e1)
	}
	sig_d, e := encoding.ReadHexScalar(curve, f1)
	if e != nil {
		fmt.Println(e)
	}

	f2, e2 := os.Open(path + "/Pubkey.txt")
	if e2 != nil {
		fmt.Println(e2)
	}
	pub_key, e_2 := encoding.ReadHexPoint(curve, f2)
	if e_2 != nil {
		fmt.Println(e_2)
	}
	path1 := "Broadcast/" + peer_number + "/Signing/KGD.txt"
	f3, e3 := os.Open(path1)
	if e3 != nil {
		fmt.Println(e3)
	}
	KGD_j, e_3 := encoding.ReadHexPoint(curve, f3)
	if e_3 != nil {
		fmt.Println(e_3)
	}

	message, e4 := ioutil.ReadFile(path + "/Message.txt")
	if e4 != nil {
		fmt.Println(e4)
	}

	newS := Signature{}
	newS.s = sig_d
	newS.r = KGD_j
	// fmt.Println("INSIDE DEComit ->> Message:", string(message))
	// fmt.Println("INSIDE DECMIT Pubkey:", pub_key)
	// fmt.Println("INSIDE DECMIT Sign:", sig_d)
	// fmt.Println("INSIDE DECMIT newS.s:", newS.s)
	t := Comit_Verify(string(message), newS, pub_key)
	//fmt.Println(t)
	if t {

		return pub_key.String()
	} else {
		return "Invalid"
	}
}
