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

	"github.com/libp2p/go-libp2p-core/protocol"

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
	fmt.Println(LagCoeff.String())
	return LagCoeff
}

func verify_R_i(Peer_Count int64, T int64) bool {
	alpha_0_sum := Get_Sum_alpha0(Peer_Count)
	sum := curve.Point().Null()
	var i int64
	for i = 1; i <= Peer_Count; i++ {
		path := "Data/" + strconv.Itoa(int(i)) + "/Signing/U_i.txt"
		file, _ := os.Open(path)
		temp, _ := encoding.ReadHexPoint(curve, file)
		lambda := Lambda(T, i)
		prod := curve.Point().Mul(lambda, temp)
		sum = sum.Add(sum, prod)
	}
	fmt.Println(alpha_0_sum.String(), "\n", sum.String())
	if alpha_0_sum.Equal(sum) {
		return true
	} else {
		return false
	}
}

func Get_Sum_alpha0(Peer_Count int64) kyber.Point {
	var i int64
	sum := curve.Point().Null()
	for i = 1; i <= Peer_Count; i++ {
		path := "Broadcast/" + strconv.Itoa(int(i)) + "/Signing/Alphas/alpha0.txt"
		file, _ := os.Open(path)
		temp, _ := encoding.ReadHexPoint(curve, file)
		sum = sum.Add(sum, temp)
	}
	return sum
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

func combine_T_Unknown(T_arr []int, peer_number, Message string) (kyber.Scalar, kyber.Point) {
	Peer_Count := len(peer_details_list)
	T := Threshold
	var Vsum kyber.Scalar = curve.Scalar().Zero()
	// err := os.MkdirAll("Received/Signing/Combine", os.ModePerm)
	// if err != nil {
	// 	panic(err)
	// }
	var Usum kyber.Point = curve.Point().Null()

	for i := 1; i <= Peer_Count; i++ {
		path := "Broadcast/" + fmt.Sprint(i) + "/Signing/V_i.txt"
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		Lambda_i := Lambda(int64(T), int64(i))
		Lambda_i2 := Lambda_i
		V_i, _ := encoding.ReadHexScalar(curve, file)

		prod := curve.Scalar().Mul(Lambda_i, V_i)
		Vsum = Vsum.Add(Vsum, prod)

		// path2 := "Broadcast/" + fmt.Sprint(i) + "/Signing/U_i.txt"
		// file, err = os.Open(path2)
		// if err != nil {
		// 	continue
		// }
		// // Lambda_i2 := Lambda(int64(T), int64(i))
		// U_i, _ := encoding.ReadHexPoint(curve, file)
		// prod2 := curve.Point().Mul(Lambda_i2, U_i)
		// Usum = Usum.Add(Usum, prod2)
		path2 := "Data/" + strconv.Itoa(int(i)) + "/Signing/U_i.txt"
		file2, _ := os.Open(path2)
		temp, _ := encoding.ReadHexPoint(curve, file2)

		prod2 := curve.Point().Mul(Lambda_i2, temp)
		Usum = Usum.Add(Usum, prod2)
	}
	fmt.Println("Sum of all V_i:", Vsum.String())
	fmt.Println("Sum of All labda U_i:", Usum.String())
	file, _ := os.OpenFile("Received/Signing/"+peer_number+"/V.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	encoding.WriteHexScalar(curve, file, Vsum)
	file, _ = os.OpenFile("Received/Signing/"+peer_number+"/U.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	encoding.WriteHexPoint(curve, file, Usum)

	GK := Get_Group_Key(int64(Peer_Count))
	// x := Vsum.Clone()
	// y := Usum.Clone()
	fmt.Println("INSIDE GKEY:", GK.String())
	res := Verify_sign_share(Vsum, Usum, Usum, Message, GK)
	if res {
		fmt.Println("SUCCESS VERIFICATION OF SIGNATURE")
	} else {
		fmt.Println("INSIDE FAILED TO VERIFIY")
	}

	return Vsum, Usum

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
func Signing_T_Unkown(U kyber.Point, x_i kyber.Scalar, Message string, peer_number string) (kyber.Scalar, kyber.Point) {

	file, _ := os.Open("Received/Signing/" + peer_number + "/R_i.txt")
	R_i, _ := encoding.ReadHexScalar(curve, file)
	U_i := curve.Point().Mul(R_i, g)

	// var T int64 = int64(Threshold)
	// j, _ := strconv.Atoi(peer_number)

	Hashing_message := Message + U.String()
	H, _ := hash_sign([]byte(Hashing_message))
	var H1 kyber.Scalar
	H1 = curve.Scalar().Zero()
	H1.SetBytes(H)
	H1 = H1.Mul(H1, x_i) //H1=H*x_i
	// H1 = H1.Mul(H1, Lambda(T, int64(j)))
	V_i := R_i.Add(R_i, H1) //Val= R_i+ H1

	return V_i, U_i
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
		// R_i := Read_Ri(peer_number)
		V_i, _ := Signing_T_Unkown(U, x_i[i], Message, peer_number)
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

func Presigning_T_Unknown(peer_number string, Peer_Count int64) {

	var protocolID protocol.ID = "/keygen/0.0.1"
	var T int64 = int64(Threshold)

	fmt.Println("******************************************PRESIGNING PHASE STARTED *******************************************")
	// var r_i kyber.Scalar
	var U_i kyber.Point
	U_i = curve.Point().Null()
	r_i := curve.Scalar().Pick(curve.RandomStream())
	U_i = U_i.Mul(r_i, g)
	// U_i, r_i := Setup_Keys(T, int64(Peer_Count), peer_number, g)
	os.MkdirAll("Data/"+peer_number+"/Signing/", os.ModePerm)
	file, _ := os.Create("Data/" + peer_number + "/Signing/r_i.txt")
	file2, _ := os.Create("Data/" + peer_number + "/Signing/U_i_temp.txt")

	encoding.WriteHexScalar(curve, file, r_i)
	encoding.WriteHexPoint(curve, file2, U_i)

	U_i_sending, _ := os.ReadFile("Data/" + peer_number + "/Signing/U_i_temp.txt")
	status_struct.Phase = 8
	send_data(peer_details_list, string(U_i_sending), "U_i_temp", protocolID)
	wait_until(8)

	//	Peer_Commitment(U_i, r_i, peer_number, Peer_Count)

	fmt.Println("Commiting Signing r_i")
	file, _ = os.Open("Data/" + peer_number + "/Signing/r_i.txt")
	// file2, _ = os.Open("Data/" + peer_number + "/Signing/U_i.txt")

	r_i, err := encoding.ReadHexScalar(curve, file)
	// U_i, err2 := encoding.ReadHexPoint(curve, file2)

	if err != nil {
		fmt.Println("error occured")
	}

	// fmt.Println("r_i->", r_i.String(), "U_i->", U_i.String())
	// path := "Data/" + peer_number + "/SSK.txt"
	// f2, _ := os.Open(path)
	// f_2, _ := encoding.ReadHexScalar(curve, f2)
	// f2.Close()
	//commiting SSK
	Commitment_sign(r_i, "helloworld", peer_number)
	//Broadcasting KGC

	// Helper.Broadcast_KGC((peer_number))

	fmt.Println("Broadcasting KGC values ....")

	fmt.Println("")
	fmt.Println("Broadcasting Signature_S ....")
	f, _ := os.ReadFile("Commitment/Signing/" + peer_number + "/KGC/Signature_S" + ".txt")
	status_struct.Phase = 9
	send_data(peer_details_list, string(f), "Signature_S", protocolID)
	// wait_until(9)

	fmt.Println("Broadcasting PubKey ....")
	f1, _ := os.ReadFile("Commitment/Signing/" + peer_number + "/KGC/PubKey" + ".txt")
	status_struct.Phase = 10
	fmt.Println("-->", string(f1))

	send_data(peer_details_list, string(f1), "PubKey", protocolID)
	// wait_until(10)

	fmt.Println("Broadcasting Message ....")
	f3, _ := os.ReadFile("Commitment/Signing/" + peer_number + "/KGC/Message" + ".txt")
	status_struct.Phase = 11
	send_data(peer_details_list, string(f3), "Message", protocolID)
	// wait_until(11)

	fmt.Println("Broadcasting KGD values ....")

	f4, _ := os.ReadFile("Commitment/Signing/" + peer_number + "/KGD" + ".txt")
	status_struct.Phase = 12
	send_data(peer_details_list, string(f4), "KGD", protocolID)
	wait_until(12)

	// case 11:
	var i int64
	//Recieving KGC from peers
	for i = 1; i <= int64(Peer_Count); i++ {
		if i == int64(my_index+1) {
			continue
		}
		Recieve_KGC_sign(strconv.Itoa(int(i)))
	}
	//Recieving KGD from peers
	for i = 1; i <= int64(Peer_Count); i++ {
		if i == int64(my_index+1) {
			continue
		}
		Recieve_KGD_sign(strconv.Itoa(int(i)))
	}

	//Decomiting Values
	for i = 1; i <= int64(Peer_Count); i++ {
		if i == int64(my_index+1) {
			continue
		}
		y_j := Decommitment_j_sign(strconv.Itoa(int(i)))
		if y_j == "Invalid" {
			fmt.Printf("Peer %s commited Wrong Values Process Aborting \n", strconv.Itoa(int(i)))
			//break
		} else {
			fmt.Printf("Peer %d Successfully Commited his values \n", i)
			fmt.Printf("Recieved Value from decommitment module is %s \n", y_j)
			fmt.Printf("\n")
		}
	}

	// case 12:
	//vss k=threshold
	// f2, _ := os.Open("Received/Signing/" + peer_number + "/G.txt")
	// x_i, _ := encoding.ReadHexScalar(curve, f2)
	// f2.Close()
	// x_i := curve.Scalar().Pick(curve.RandomStream())

	//Set_secret_sign(x_i)

	poly := []kyber.Scalar{}  // to store coefficients
	share := []kyber.Scalar{} // to store share
	alphas := []kyber.Point{} // to store alphas

	// pt:=curve.Scalar()
	// pt.SetBytes()
	// var i int64

	for i = 0; i < T; i++ {
		poly = append(poly, curve.Scalar().Zero())
	}

	for i = 0; i < T; i++ {
		alphas = append(alphas, curve.Point().Null())
	}

	for i = 1; i <= int64(Peer_Count); i++ {
		share = append(share, curve.Scalar().Zero())
	}

	// to generate coefficients of the polynomial         //r_i
	Generate_Polynomial_coefficients(T, poly, peer_number, r_i, "vss/Signing/"+peer_number)
	// fmt.Println("COFFE", poly[0].String(), "\n", poly[1].String(), "\n")

	Generate_share(int64(Peer_Count), T, poly, share, peer_number, "vss/Signing/"+peer_number)
	// fmt.Println("SHARES", share[0].String(), "\n", share[1].String(), "\n")

	//Generating Alphas
	Generate_Alphas(T, alphas, poly, peer_number, "vss/Signing/"+peer_number)
	// fmt.Println("ALPHAS", alphas[0].String(), "\n", alphas[1].String(), "\n")

	//Broadcasting alphas

	status_struct.Phase = 13
	for i = 0; i < T; i++ {
		send_data(peer_details_list, alphas[i].String(), fmt.Sprint(i), protocolID)
	}

	wait_until(13)
	Recieve_Alphas_sign(int64(Peer_Count), peer_number, T)

	//Broadcasting Share
	status_struct.Phase = 14
	paths := "vss/Signing/" + peer_number
	for i = 1; i <= int64(Peer_Count); i++ {
		_f, _ := os.Open(paths + "/Indivisual_Share" + strconv.Itoa(int(i)) + ".txt")
		shares, _ := encoding.ReadHexScalar(curve, _f)
		tosend := shares.String()
		fmt.Println("TO SeND:", tosend)
		send_data(peer_details_list, tosend, fmt.Sprint(i), protocolID)
	}
	wait_until(14)

	//Receiving Sign shares
	Recieve_Share_sign(peer_number, int64(Peer_Count))
	// var p int
	// for p = 0; p <= Peer_Count; p++ {
	// 	if strconv.Itoa(p) == peer_number {
	// 		continue
	// 	}
	// 	alphas2 := []kyber.Point{} // to store alphas
	// 	for i = 0; i < T; i++ {
	// 		alphas2 = append(alphas2, curve.Point().Null())
	// 	}

	// 	var q int
	// 	for q = 0; q < int(T); q++ {
	// 		fx, err := os.Open("vss/Signing/" + strconv.Itoa(int(p)) + "/alpha" + strconv.Itoa(int(q)) + ".txt")
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fx_p, _ := encoding.ReadHexPoint(curve, fx)
	// 		alphas2[q] = fx_p
	// 	}
	// 	fx2, errx := os.Open("vss/Signing/" + strconv.Itoa(int(p)) + "/Indivisual_Share" + peer_number + ".txt")
	// 	if errx != nil {
	// 		fmt.Print("ERROR ERRX")
	// 	}
	// 	my_share, _ := encoding.ReadHexScalar(curve, fx2)
	// 	check(peer_number, my_share, T, alphas2)

	// }

	// case 13:
	fmt.Println("Verifying Signing Shares")

	path := "Received/Signing/" + peer_number + "/R_i.txt"
	R_i := Verify_Share(peer_number, int64(Peer_Count), T, true)
	file, _ = os.Create(path)
	encoding.WriteHexScalar(curve, file, R_i)

	U_i = curve.Point().Mul(R_i, g)
	file, _ = os.Create("Data/" + peer_number + "/Signing/U_i.txt")
	encoding.WriteHexPoint(curve, file, U_i)

	if verify_R_i(int64(Peer_Count), T) {
		fmt.Println("VERIFIED Ri")
	} else {
		fmt.Println("NOT VERIFIED Ri")
	}

	// U_i_sending, _ := os.ReadFile("Data/" + peer_number + "/Signing/U_i.txt")

	// status_struct.Phase = 8
	// send_data(peer_details_list, string(U_i_sending), "U_i", protocolID)
	// wait_until(8)

	file3, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	encoding.WriteHexScalar(curve, file3, R_i)
	file3.Close()

	//TAKE IT ALL INTO COMBINE FCN
	// // case 14:
	// var U kyber.Point
	// U = curve.Point().Null()
	// fmt.Println("U_i's : ")
	// // var i int
	// for i = 1; i <= int64(Peer_Count); i++ {
	// 	if i == int64(my_index+1) {
	// 		path := "Data/" + fmt.Sprint(i) + "/Signing/U_i.txt"
	// 		file, err := os.Open(path)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		U_i, e1 := encoding.ReadHexPoint(curve, file)
	// 		if e1 != nil {
	// 			panic(e1)
	// 		}
	// 		fmt.Println(U_i)
	// 		if U.Equal(curve.Point().Null()) {
	// 			U = U_i
	// 		} else {
	// 			U.Add(U, U_i)
	// 		}
	// 	} else {
	// 		// peer := strconv.Itoa(int(i))
	// 		// U_i := Read_Ui(peer)
	// 		path := "Data/" + fmt.Sprint(i) + "/Signing/U_i.txt"
	// 		file, err := os.Open(path)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		U_i, e1 := encoding.ReadHexPoint(curve, file)
	// 		if e1 != nil {
	// 			panic(e1)
	// 		}
	// 		fmt.Println(U_i)

	// 		if U.Equal(curve.Point().Null()) {
	// 			U = U_i
	// 		} else {
	// 			U.Add(U, U_i)
	// 		}
	// 	}
	// }
	// // //U,U_i,R_i is generated for each peer above
	// fmt.Println("U:")
	// fmt.Println(U.String())

	U := Get_Sum_alpha0(Peer_Count)
	os.MkdirAll("Data/"+peer_number+"/Signing/", os.ModePerm)
	file, _ = os.Create("Data/" + peer_number + "/Signing/U.txt")
	encoding.WriteHexPoint(curve, file, U)
	// choice = 3
	// time.Sleep(time.Second * 2)

}

func Signing(peer_number, Message string) {

	fmt.Println("MESSAGE TO SIGN:", Message)
	var protocolID protocol.ID = "/keygen/0.0.1"
	Peer_Count := len(peer_details_list)
	T := Threshold
	// r_i:= curve.Scalar().Pick(curve.RandomStream())

	fmt.Printf("********************************************* SIGNING PHASES STARTED ******************************************\n")

	file, _ := os.Open("Received/Signing/" + peer_number + "/R_i.txt")
	// R_i, _ := encoding.ReadHexScalar(curve, file)
	// U_i := curve.Point().Mul(R_i, g)

	file, _ = os.Open("Received/" + peer_number + "/G.txt")
	x_i, _ := encoding.ReadHexScalar(curve, file)

	file, _ = os.Open("Data/" + peer_number + "/Signing/U.txt")
	U, _ := encoding.ReadHexPoint(curve, file)
	fmt.Println("U from PreSign:", U.String())

	// file, _ = os.Open("Data/" + peer_number + "/Signing/U_i.txt")
	// U_i, _ := encoding.ReadHexPoint(curve, file)

	// fmt.Println("U_i READ:", U_i.String(), "\n")
	// file, _ = os.Open("Data/" + peer_number + "/Signing/U.txt")
	// U, _ := encoding.ReadHexPoint(curve, file)
	V_i, U_i := Signing_T_Unkown(U, x_i, Message, peer_number)
	fmt.Println("U_i returned from sign:", U_i.String(), "\n")

	X_i := curve.Point().Mul(x_i, g)
	fmt.Println("X_i", X_i.String())

	// z1 := curve.Scalar().Pick(curve.RandomStream())
	// z2 := curve.Scalar().Pick(curve.RandomStream())
	// X_i := curve.Point().Mul(z2, g)

	// Z1 := curve.Point().Mul(z1, g)

	// Hashing_message := Message + U.String()
	// h, _ := hash_sign([]byte(Hashing_message))
	// var H1 kyber.Scalar
	// H1 = curve.Scalar().Zero()
	// H1.SetBytes(h)

	// prod := curve.Scalar().Mul(H1, z2)
	// Z3 := curve.Scalar().Add(z1, prod)

	// if Verify_sign_share(Z3, U, Z1, Message, X_i) {
	// 	fmt.Println("Test SHARES ARE VERIFIED")
	// } else {
	// 	fmt.Println("NOT VERIFIED TEst Shares")
	// }

	if Verify_sign_share(V_i, U, U_i, Message, X_i) {
		fmt.Println("INDIVIDUAL SHARES ARE VERIFIED")
	} else {
		fmt.Println("NOT VERIFIED INDIVIDUAL SHARES")
	}

	file, _ = os.Create("Data/" + peer_number + "/Signing/V_i.txt")
	encoding.WriteHexScalar(curve, file, V_i)

	//Broadcasting V_i
	status_struct.Phase = 15

	fmt.Println(this_vault, my_index, peer_details_list)
	tosend, _ := encoding.ScalarToStringHex(curve, V_i)
	send_data(peer_details_list, tosend, "V_i", protocolID)

	Wait_until_for_sign(15, T)

	// status_struct.Phase = 16

	// tosend, _ = encoding.PointToStringHex(curve, U)
	// send_data(peer_details_list, tosend, "U", protocolID)

	// wait_until(16)
	// choice = 4
	// time.Sleep(time.Second * 2)
	fmt.Println("************ COMBINATION PHASE ****************")
	T_arr := [...]int{1} //not actually used
	fmt.Println(Peer_Count)
	GKey := Get_Group_Key(int64(Peer_Count))
	Vsum, Usum := combine_T_Unknown(T_arr[:], peer_number, Message)
	fmt.Println("************ VERIFYING ****************")
	// file, _ = os.Open("Received/Signing/" + peer_number + "/V.txt")
	// V, _ := encoding.ReadHexScalar(curve, file)
	// file, _ = os.Open("Received/Signing/" + peer_number + "/U.txt")
	// U2, _ := encoding.ReadHexPoint(curve, file)
	fmt.Println("GKEY:", GKey.String())
	fmt.Println("VSUM:", Vsum.String())
	fmt.Println("Usum:", Usum.String())

	res := Verify_sign_share(Vsum, Usum, Usum, Message, GKey)
	if res {
		fmt.Println("INSIDe SUCCESS VERIFICATION OF SIGNATURE")
	} else {
		fmt.Println("INSIDE FAILED TO VERIFIY")
	}

}
