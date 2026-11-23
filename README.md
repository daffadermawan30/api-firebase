# Firebase Authentication (api_dart)

## Anggota & Pembagian Tugas
- Nama : Daffa Dermawan
- NIM : 1124160177

Aplikasi Flutter untuk autentikasi login (Email & Password) yang terintegrasi Firebase.

## Langkah Integrasi Firebase
1. **Firebase Console**: Buat project Firebase (`project-ondeva`) dan aktifkan menu **Authentication > Sign-in method > Email/Password**.
2. **FlutterFire CLI**: Jalankan `flutterfire configure --project=project-ondeva` untuk generate `firebase_options.dart`.
3. **Dependencies**: Tambahkan `firebase_core` dan `firebase_auth` di `pubspec.yaml`, lalu jalankan `flutter pub get`.
4. **Inisialisasi**: Panggil `WidgetsFlutterBinding.ensureInitialized()` dan `await Firebase.initializeApp(...)` di `main()`.
5. **Login Auth**: Gunakan fungsi `FirebaseAuth.instance.signInWithEmailAndPassword()` pada form login.

## Cara Menjalankan
```bash
flutter pub get
flutter run
```
