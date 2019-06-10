'use strict';

const BLOCK_SIZE = 16;

const pkcs7 = require('pkcs7'),
  aes = require('aes-js'),
  hash = require('hash.js'),
  salt = new Uint8Array([
    0x3f, 0xb8, 0x87, 0x7d, 0x37, 0xfd, 0xc0, 0x4e,
    0x4a, 0x47, 0x65, 0xEF, 0xb8, 0xab, 0x7d, 0x36,
  ]);

function str_to_buf(s) {
  if (typeof (s) == 'string') {
    const len = s.length;
    const buf = new Uint8Array(len);
    for (let i = 0; i != len; ++i) {
      buf[i] = s.charCodeAt(i);
    }
    return buf;
  } else {
    return s;
  }
}

function key_derivation(passwd) {
  const h = hash.sha256().update(hash.sha256().update(passwd).update(salt).digest()).digest();
  return new Uint8Array(h.slice(0, BLOCK_SIZE));
}

function rand_bytes(sz) {
  const iv = new Uint8Array(sz);
  for (let i = 0; i != sz; ++i) {
    iv[i] = Math.floor(Math.random() * 256);
  }
  return iv;
}

function pad(s) {
  return pkcs7.pad(str_to_buf(s));
}

function unpad(s) {
  if (!s || s.length < BLOCK_SIZE) {
    throw new Error('invalid pkcs#7 data');
  }

  const p = s[s.length - 1];
  if (p > BLOCK_SIZE) {
    throw new Error('invalid pkcs#7 padding length');
  }

  return pkcs7.unpad(s);
}

function encrypt(s, passwd) {
  const p = key_derivation(passwd);
  const iv = rand_bytes(BLOCK_SIZE);
  const cbc = new aes.ModeOfOperation.cbc(p, iv);
  const pad_in = pad(s);
  const ret = new Uint8Array(iv.length + pad_in.length);
  const ed = cbc.encrypt(pad_in);
  if (!ed || !ed.length) {
    throw new Error('encrypt failed');
  }
  ret.set(iv);
  ret.set(ed, iv.length);
  return ret;
}

function decrypt(c, passwd) {
  const p = key_derivation(passwd);
  const cb = str_to_buf(c);

  if (cb.length < 2 * BLOCK_SIZE) {
    throw new Error('invalid cipher');
  }

  const iv = cb.slice(0, BLOCK_SIZE);
  const cbc = new aes.ModeOfOperation.cbc(p, iv);
  const dd = cbc.decrypt(cb.slice(BLOCK_SIZE));

  if (!dd || !dd.length || dd.length != (cb.length - BLOCK_SIZE)) {
    throw new Error('invalid decrypted data');
  }

  return unpad(new Uint8Array(dd));
}

function decrypt_string(c, passwd) {
  const d = decrypt(c, passwd);
  return aes.utils.utf8.fromBytes(d);
}

const e2e = {
  encrypt: encrypt,
  decrypt: decrypt,
  decrypt_string: decrypt_string,
};

export default e2e;
