package main

import (
	SHA_256 "crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	// cmndcm "keygen/commitment_decommitment"
	// "keygen/vss"
	"math/big"

	"os"
	"strconv"

	// "keygen/Helper"

	"gopkg.in/dedis/kyber.v2"
	// "gopkg.in/dedis/kyber.v2/group/edwards25519"
	"gopkg.in/dedis/kyber.v2/util/encoding"
)

// var secret big.Int

// var curve = edwards25519.NewBlakeSHA256Ed25519()

func Lambda(t, j int64) kyber.Scalar {
	var i int64
	den := curve.Scalar().One()
	var LagCoeff = curve.Scalar().One()        //
	var J kyber.Scalar = curve.Scalar().Zero() //Converting j to kyber scalar from int64
	J.SetInt64(j)
	for i = 1; i <= t; i++ {
		if i == j {
			continue
		}
		var I kyber.Scalar = curve.Scalar().Zero()
		I.SetInt64(i)
		den.Sub(I, J)               //den=(i-j)
		den.Inv(den)                //1/(i-j)
		den.Mul(den, I)             //i/(i-j)
		LagCoeff.Mul(LagCoeff, den) // product (i/(i-j)) for each i from 1 to t such that i!=j
	}
	return LagCoeff
}

//Hash the given byte array using SHA256
func hash_sign(value []byte) ([]byte, error) {
	h := SHA_256.New()
	h.Write(value)
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	ret, _ := hex.DecodeString(sha1_hash)
	return ret, nil
}

func Broadcast_verification_set(peer_Count, Peer_number int64, share []kyber.Scalar) {
	// fmt.Printf("share is %d", share[1])
	path := "Broadcast/VerificationSet/SetBy" + fmt.Sprint(Peer_number)
	os.MkdirAll(path, os.ModePerm)
	for i := 1; i <= int(peer_Count); i++ {
		path1 := path + "/F_" + fmt.Sprint(Peer_number) + "(" + fmt.Sprint(i) + ")" + ".txt"
		file, err := os.OpenFile(path1, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			panic(err)
		}
		encoding.WriteHexScalar(curve, file, share[i])
		file.Close()
	}
}

func Broadcast_sigma_T_Unknown(U, U_i kyber.Point, V_i kyber.Scalar, peer_number int64, Message string) {
	path := "Broadcast/Sigmas_T_Unknown/sigma" + fmt.Sprint(peer_number)
	os.MkdirAll(path, os.ModePerm)
	file, _ := os.OpenFile(path+"/U.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	fmt.Fprint(file, U.String())
	file.Close()
	file, _ = os.OpenFile(path+"/U"+fmt.Sprint(peer_number)+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	fmt.Fprint(file, U_i.String())
	file.Close()
	file, _ = os.OpenFile(path+"/V"+fmt.Sprint(peer_number)+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	fmt.Fprint(file, V_i.String())
	file.Close()
	file, _ = os.OpenFile(path+"/Message"+fmt.Sprint(peer_number)+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	fmt.Fprint(file, Message)
	file.Close()
}

func Broadcast_sigma_T_known(U, U_i kyber.Point, V_i *big.Float, peer_number int64) {
	path := "Broadcast/Sigmas_T_Known/sigma" + fmt.Sprint(peer_number)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(path+"/U.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(file, U.String())
	file, err = os.OpenFile(path+"/U"+fmt.Sprint(peer_number)+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(file, U_i.String())
	file, err = os.OpenFile(path+"/V"+fmt.Sprint(peer_number)+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(file, V_i.String())
}

func Broadcast_Ui(peer_number int64, U_i kyber.Point) {
	path := "Broadcast/U"
	os.MkdirAll(path, os.ModePerm)
	path = path + "/U_" + fmt.Sprint(peer_number)
	f1, err := os.OpenFile(path+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	encoding.WriteHexPoint(curve, f1, U_i)
	f1.Close()
}

func combine_T_Unknown(T_arr []int, peer_number string) {
	l := len(T_arr)
	var sum kyber.Scalar = curve.Scalar().Zero()
	// err := os.MkdirAll("Received/Signing/Combine", os.ModePerm)
	// if err != nil {
	// 	panic(err)
	// }
	for i := 0; i < l; i++ {
		path := "Broadcast/" + fmt.Sprint(i) + "/Signing/V_i.txt"
		file, _ := os.Open(path)
		Lambda_i := Lambda(int64(l), int64(i))
		V_i, _ := encoding.ReadHexScalar(curve, file)
		prod := Lambda_i.Mul(Lambda_i, V_i)
		sum = sum.Add(sum, prod)
	}
	fmt.Println("Sum of all V_i:", sum.String())
	file, _ := os.OpenFile("Received/Signing/"+peer_number+"/V.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	fmt.Fprint(file, sum.String())
}

func combine_T_known(T_arr []int) {
	l := len(T_arr)
	sum := big.NewFloat(0)
	os.MkdirAll("Private/Combine", os.ModePerm)
	for i := 0; i < l; i++ {
		path := "Broadcast/Sigmas_T_Known/sigma" + fmt.Sprint(T_arr[i]) + "/V" + fmt.Sprint(T_arr[i]) + ".txt"
		data, _ := ioutil.ReadFile(path)
		//temp := new(big.Int)
		//V_i, _ := temp.SetString(string(data), 10)

		V_i_float := new(big.Float)
		V_i_float, _ = V_i_float.SetString(string(data))

		//prod := new(big.Float)

		sum = sum.Add(sum, V_i_float)
	}
	file, _ := os.OpenFile("Private/Combine/V.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	fmt.Fprint(file, sum.String())

}

func Signing_T_Known(U kyber.Point, x_i kyber.Scalar, r_i kyber.Scalar, Message string) kyber.Scalar {
	Hashing_message := U.String() + Message
	H, _ := hash_sign([]byte(Hashing_message))
	var H1 kyber.Scalar
	H1 = curve.Scalar().Zero()
	H1.SetBytes(H)
	H1 = H1.Mul(H1, x_i)    //H1=H*x_i
	V_i := r_i.Add(r_i, H1) //Val= r_i+ H1
	return V_i
}

func T_Known_Setup(Peer_Count int64, peer_number string) {
	path := "Private/Delta" + peer_number
	r_i := curve.Scalar().Pick(curve.RandomStream())
	Store_Scaler(r_i, path, "r_i")
	g := curve.Point().Base()

	U_i := g.Mul(r_i, g)

	//commiting r_i
	message := U_i.String() //message =U_i ,p=r_i
	fmt.Printf("message hai %s \n", message)
	Commitment(r_i, message, path)
	//Broadcasting KGC
	Broadcast_KGC(peer_number)
	//Broadcasting KGD
	Broadcast_KGD(peer_number)
	//Brodcasting the value of u_i
	peer_num_64, _ := strconv.ParseInt(peer_number, 10, 64)
	Broadcast_Ui(peer_num_64, U_i)
}

func Sign_T_Known(K int64, Peer_Count int64, Message string) {

	//---------------PreSigning-------------
	var i int64
	for i = 1; i <= Peer_Count; i++ {
		peer_number := strconv.Itoa(int(i))
		T_Known_Setup(Peer_Count, peer_number)
	}
	Receiving_Commitment(Peer_Count)
	Decommitment(Peer_Count)
	var U kyber.Point = curve.Point().Null()
	for i = 1; i <= Peer_Count; i++ {
		peer_number := strconv.Itoa(int(i))
		U_i := Read_Ui(peer_number)
		U.Add(U, U_i)
	}
	// /*---------------Signing----------------*/
	// for i = 1; i <= Peer_Count; i++ {
	// 	peer_number := strconv.Itoa(int(i))
	// 	Signing_T_Known(K, Peer_Count, peer_number, Message, U)

	// }
	// //_______________Combining--------------
	// T_arr := [...]int{1}
	// combine_T_known(T_arr[:])

}

/********************Signing When T unknown ****************************/

//Generates Random Secret Number and Calculates U_i and broadcasts it
func Setup_Keys(T int64, Peer_Count int64, peer_number string, g kyber.Point) (kyber.Point, kyber.Scalar) {
	//Getting Random Secret Number
	r_i := curve.Scalar().Pick(curve.RandomStream())
	//Storing r_i into private folder
	// path := "Data/" + peer_number + "/"
	// file, err := os.OpenFile(path+"r_i.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	// if err != nil {
	// 	panic(err)
	// }
	// encoding.WriteHexScalar(curve, file, r_i)
	// file.Close()
	//U_i=g*r_i
	U_i := g.Mul(r_i, g)
	//Brodcasting the value of u_i
	return U_i, r_i
}

func Peer_Commitment(U_i kyber.Point, r_i kyber.Scalar, peer_number string, Peer_Count int64) {
	Message := U_i.String() //message =U_i ,p=r_i
	fmt.Printf("Message  %s \n", Message)
	Commitment(r_i, Message, peer_number)
	//Broadcasting KGC
	Broadcast_KGC(peer_number)
	//Broadcasting KGD
	Broadcast_KGD(peer_number)
}

func Receiving_Commitment(Peer_Count int64) {
	var i int64
	//Recieving KGC from peers
	for i = 1; i <= int64(Peer_Count); i++ {
		Recieve_KGC(strconv.Itoa(int(i)))
	}
	//Recieving KGD from peers
	for i = 1; i <= int64(Peer_Count); i++ {
		Recieve_KGD(strconv.Itoa(int(i)))
	}
}

func Decommitment(N int64) {
	var i int64
	fmt.Println(" \n ********************Running Decommitment***************** \n ")
	for i = 1; i <= int64(N); i++ {
		y_j := Decommitment_j(strconv.Itoa(int(i)))
		if y_j == "Invalid" {
			fmt.Printf("Peer %s commited Wrong Values Process Aborting \n", strconv.Itoa(int(i)))
			//break
		} else {
			fmt.Printf("Peer %d Successfully Commited his values \n", i)
			fmt.Printf("Recieved Value from decommitment module is %s \n", y_j)
			fmt.Printf("\n")
		}
	}
}

// func Peer_VSS_Setup(T int64, N int64, peer_number string, x_i kyber.Scalar) {
// 	var i int64
// 	//Setting Secret
// 	Set_secret(x_i)
// 	poly := []kyber.Scalar{}  // to store coefficients
// 	share := []kyber.Scalar{} // to store share
// 	alphas := []kyber.Point{} // to store alphas

// 	for i = 0; i <= T; i++ {
// 		poly = append(poly, curve.Scalar().Zero())
// 	}

// 	for i = 0; i <= T; i++ {
// 		alphas = append(alphas, curve.Point().Null())
// 	}

// 	for i = 0; i <= int64(N); i++ {
// 		share = append(share, curve.Scalar().Zero())
// 	}
// 	//Generating Polynomial coefficients
// 	Generate_Polynomial_coefficients(T, poly, peer_number)
// 	//Generating Shares
// 	Generate_share(N, T, poly, share, peer_number)
// 	//Generating Alphas
// 	Generate_Alphas(T, alphas, poly, peer_number)
// 	//Broadcasting alphas
// 	//--->>Brodcast_Alphas(peer_number, alphas, int64(T))
// 	fmt.Printf("\n")
// }

//verify_each_share
// func verify_each_share(peer_number string, peer_count int64, share []kyber.Scalar, T int64) {
// 	var i int64
// 	fmt.Printf("Verifying shares Recieved to %s \n", peer_number)
// 	for i = 1; i <= peer_count; i++ {
// 		var j int64
// 		path := "Recieved/ALPHAS/Alpha" + strconv.Itoa(int(i)) + "/"
// 		alphas := []kyber.Point{} // to store alphas

// 		for j = 0; j <= T; j++ {
// 			alphas = append(alphas, curve.Point().Null())
// 		}
// 		for j = 0; j < T; j++ {
// 			f1, e1 := os.Open(path + "alpha" + strconv.Itoa(int(j)) + ".txt")
// 			if e1 != nil {
// 				panic(e1)
// 			}
// 			alpha, _ := encoding.ReadHexPoint(curve, f1)
// 			alphas[j] = alpha
// 			f1.Close()
// 		}
// 		//fmt.Println(alphas)
// 		I, err := strconv.Atoi(peer_number)
// 		if err != nil {
// 			panic(err)
// 		}
// 		if !vss.Verify_i(int64(I), share[i], T, alphas) {
// 			fmt.Printf("Peer %d shared wrong values mission aborting \n", i)
// 		} else {
// 			fmt.Printf("Shared Verified for Peer %d \n", i)
// 		}
// 	}
// }

// func Verify_Share(peer_number string, N int64, T int64) kyber.Scalar {
// 	var i int64
// 	path := "Recieved/sharesfor" + peer_number + "/share"
// 	share := []kyber.Scalar{} // to store share
// 	for i = 0; i <= int64(N); i++ {
// 		share = append(share, curve.Scalar().Zero())
// 	}
// 	for i = 1; i <= N; i++ {
// 		file, err := os.Open(path + strconv.Itoa(int(i)) + ".txt")
// 		if err != nil {
// 			panic(err)
// 		}
// 		val, e1 := encoding.ReadHexScalar(curve, file)
// 		if e1 != nil {
// 			panic(e1)
// 		}
// 		share[i] = val
// 	}
// 	verify_each_share(peer_number, N, share, T)
// 	var R_i kyber.Scalar = curve.Scalar().Zero()
// 	for i = 1; i <= int64(N); i++ {
// 		R_i.Add(share[i], R_i)
// 	}
// 	return R_i
// }

//Function used for signing
func Signing_T_Unkown(U kyber.Point, x_i kyber.Scalar, R_i kyber.Scalar, Message string) kyber.Scalar {
	Hashing_message := U.String() + Message
	H, _ := hash_sign([]byte(Hashing_message))
	var H1 kyber.Scalar
	H1 = curve.Scalar().Zero()
	H1.SetBytes(H)
	H1 = H1.Mul(H1, x_i)    //H1=H*x_i
	V_i := R_i.Add(R_i, H1) //Val= R_i+ H1
	return V_i
}

func Read_Ui(peer_number string) kyber.Point {
	path := "Broadcast/" + peer_number + "/Signing/U_i.txt"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	U_i, e1 := encoding.ReadHexPoint(curve, file)
	if e1 != nil {
		panic(e1)
	}
	return U_i
}
func Read_Vi(peer_number string) kyber.Scalar {
	path := "Broadcast/V_" + peer_number + ".txt"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	V_i, _ := encoding.ReadHexScalar(curve, file)
	file.Close()
	return V_i
}
func Read_Ri(peer_number string) kyber.Scalar {
	path := "Private/" + "R_" + peer_number + ".txt"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	R_i, _ := encoding.ReadHexScalar(curve, file)
	file.Close()
	return R_i
}

//signing scheme when the threshold set is not known
func Sign_T_Unknown(T int64, Peer_Count int64, g kyber.Point, Message string, x_i []kyber.Scalar) { //K=threshold ,peer_Count=total number of peers in the network
	/*************************Presigning ********************/
	fmt.Println("**************Signing Module Starts Here *****************")
	var i int64
	var r_i kyber.Scalar
	var U_i kyber.Point
	for i = 1; i <= Peer_Count; i++ {
		peer_number := strconv.Itoa(int(i))
		U_i, r_i = Setup_Keys(T, Peer_Count, peer_number, g)
		Peer_Commitment(U_i, r_i, peer_number, Peer_Count)
	}
	Receiving_Commitment(Peer_Count)
	Decommitment(Peer_Count)
	for i = 1; i <= Peer_Count; i++ {
		peer_number := strconv.Itoa(int(i))
		path := "Private/" + "R_" + peer_number + ".txt"
		R_i := Verify_Share(peer_number, Peer_Count, T, true)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			panic(err)
		}
		encoding.WriteHexScalar(curve, file, R_i)
		file.Close()
	}
	var U kyber.Point = curve.Point().Null()
	fmt.Println("U_i's : ")
	for i = 1; i <= Peer_Count; i++ {
		peer_number := strconv.Itoa(int(i))
		U_i := Read_Ui(peer_number)
		fmt.Println(U_i)
		U.Add(U, U_i)
	}
	// //U,U_i,R_i is generated for each peer above
	fmt.Println("U:")
	fmt.Println(U)
	/*******************************Signing**********************/
	for i = 1; i <= Peer_Count; i++ {
		peer_number := strconv.Itoa(int(i))
		R_i := Read_Ri(peer_number)
		V_i := Signing_T_Unkown(U, x_i[i], R_i, Message)
		path := "Broadcast/V_" + peer_number + ".txt"
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			panic(err)
		}
		encoding.WriteHexScalar(curve, file, V_i)
		file.Close()
	}
	fmt.Println("V_i's: ")
	for i = 1; i <= Peer_Count; i++ {
		peer_number := strconv.Itoa(int(i))
		V_i := Read_Vi(peer_number)
		fmt.Println(V_i)
	}
	//Broadcasting Sigma Set
	for i = 1; i <= Peer_Count; i++ {
		peer_number := strconv.Itoa(int(i))
		U_i := Read_Ui(peer_number)
		V_i := Read_Vi(peer_number)
		Broadcast_sigma_T_Unknown(U, U_i, V_i, i, Message)
	}
	/********************************Combine********************/
	//T_arr := [...]int{1}
	//combine_T_Unknown(T_arr[:], "")
}
