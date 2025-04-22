package mmc

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/test"
)

func (s *Service) SelfTest() {
	s.testSecretKeyConfigs()
	s.testRSA()
	s.testFingerprint()
}

func (s *Service) testSecretKeyConfigs() {
	if _, err := s.gcmKey(); err != nil {
		panicMsg("testSecretKeyConfigs GCM key failed: %s", err.Error())
	}
}

func (s *Service) testRSA() {
	rsaKeyPairs := [][2]func() (string, error){
		//{s.UserPasswordRSAPubkeyDERBase64, s.userPasswordRSAPrivkeyDER},
		{s.rsaPubDERBase64, s.rsaPrivDER},
	}
	panicE("testRSA", test.TestRSA(rsaKeyPairs))
}

// record=`base64(rsaDeskey) base64(desRecord)`
// 第一部分为 RSA 加密后的deskey，deskey 必须很短，目前只能是8字节
// 第二部分为通过deskey加密后的数据
func (s *Service) testGenerateRecord() ([]byte, string, *ae.Error) {
	deskey := coding.RandAlphabets(8, 0)
	pubkeyDERB64, _ := s.rsaPubDERBase64()
	pem := coding.RasToPKCS8([]byte(pubkeyDERB64), false, true)
	record := generateFingerprintRecord()
	// 客户端使用临时key，可以使用EDB模式
	recordCipher, err := coding.EcbEncryptToBase64(record, []byte(deskey))
	if err != nil {
		return nil, "", ae.NewE("self-test generate fingerprint record failed encrypt record: " + err.Error())
	}
	deskeyCipherBase64, err := coding.RsaEncryptToBase64([]byte(deskey), pem, false)
	if err != nil {
		return nil, "", ae.NewE("self-test generate fingerprint record failed encrypt deskey: " + err.Error())
	}
	recordCipherBase64 := base64.StdEncoding.EncodeToString(recordCipher)
	if deskeyCipherBase64 == "" || recordCipherBase64 == "" {
		return nil, "", ae.NewE("self-test generate fingerprint record failed encrypt deskey: (%s) (%s)", deskeyCipherBase64, recordCipherBase64)
	}
	r := deskeyCipherBase64 + " " + recordCipherBase64
	return record, r, nil
}

func (s *Service) testFingerprint() {
	wantApollo := "APOLLO"
	wantUserAgent := "Fake User-Agent Aario"
	wantIP := "127.0.0.1"
	wantUUID := "123456789-0"
	wantRecord, encryptedRecord, e := s.testGenerateRecord()
	if e != nil {
		panicE("testFingerprint", e)
	}
	fp, e := s.EncryptClientRecordToFingerprint(context.Background(), []byte(encryptedRecord), "APOLLO", "Fake User-Agent Aario", "127.0.0.1")
	if e != nil {
		panicE("testFingerprint", e)
	}
	record, _, apollo, userAgent, ip, uuid, e := s.DecryptFingerprint(fp)
	if e != nil {
		panicMsg("testFingerprint decrypt fingerprint %s", e.Text())
	}
	if bytes.Compare(record, wantRecord) != 0 || wantApollo != apollo || wantUserAgent != userAgent || wantIP != ip || wantUUID != uuid {
		panicMsg("testFingerprint user auth rsa failed: decrypt fingerprint not match %s %s\n%s %s\n%s %s\n%s %s\n%s %s", string(wantRecord), string(record), wantApollo, apollo, wantUserAgent, userAgent, wantIP, ip, wantUUID, uuid)
	}
	return
}

func generateFingerprintRecord() []byte {
	header := fingerprintHeader{
		X:        0,
		Y:        0,
		Width:    0,
		Height:   0,
		XA:       0,
		YA:       0,
		WidthA:   0,
		HeightA:  0,
		ClientMs: 0,
		ServerMs: 0,
	}
	body1 := fingerprintBodySegment{
		Tag:                []byte{'k'},
		X1:                 0,
		Y1:                 0,
		XA:                 0,
		YA:                 0,
		ClientTimeOffsetMs: 0,
	}
	body2 := fingerprintBodySegment{
		Tag:                []byte{'k'},
		X1:                 0,
		Y1:                 0,
		XA:                 0,
		YA:                 0,
		ClientTimeOffsetMs: 0,
	}
	headers := header.serialize()
	bodies1 := body1.serialize()
	bodies2 := body2.serialize()

	record := make([]byte, 0, len(headers)+len(bodies1)+len(bodies2)+2)
	record = append(record, headers...)
	record = append(record, '\n')
	record = append(record, bodies1...)
	record = append(record, '\n')
	record = append(record, bodies2...)
	return record
}
