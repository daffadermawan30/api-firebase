import 'package:flutter/material.dart';
import 'package:firebase_auth/firebase_auth.dart';

class UserInfoScreen extends StatefulWidget {
  final User user;

  const UserInfoScreen({super.key, required this.user});

  @override
  State<UserInfoScreen> createState() => _UserInfoScreenState();
}

class _UserInfoScreenState extends State<UserInfoScreen> {
  late User _currentUser;

  @override
  void initState() {
    super.initState();
    _currentUser = widget.user;
    _refreshUser();
  }

  @override
  void didUpdateWidget(covariant UserInfoScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.user.uid != widget.user.uid) {
      _currentUser = widget.user;
      _refreshUser();
    }
  }

  Future<void> _refreshUser() async {
    try {
      await FirebaseAuth.instance.currentUser?.reload();
      final updatedUser = FirebaseAuth.instance.currentUser;
      if (mounted && updatedUser != null) {
        setState(() {
          _currentUser = updatedUser;
        });
      }
    } catch (_) {
      // Tetap gunakan data user yang tersimpan jika gagal memuat ulang
    }
  }

  Future<void> _logout() async {
    await FirebaseAuth.instance.signOut();
  }

  Widget _buildAvatar() {
    final photoUrl = _currentUser.photoURL;
    if (photoUrl == null || photoUrl.isEmpty) {
      return const Icon(Icons.account_circle, size: 100, color: Colors.blue);
    }

    return CircleAvatar(
      radius: 50,
      backgroundColor: Colors.blue.shade50,
      child: ClipOval(
        child: Image.network(
          photoUrl,
          width: 100,
          height: 100,
          fit: BoxFit.cover,
          loadingBuilder: (context, child, loadingProgress) {
            if (loadingProgress == null) return child;
            return const CircularProgressIndicator(strokeWidth: 2);
          },
          errorBuilder: (context, error, stackTrace) {
            return const Icon(
              Icons.account_circle,
              size: 100,
              color: Colors.blue,
            );
          },
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Dashboard User'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _refreshUser,
            tooltip: 'Refresh Profil',
          ),
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: _logout,
            tooltip: 'Logout',
          ),
        ],
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              _buildAvatar(),
              const SizedBox(height: 24),
              const Text(
                'Login Berhasil!',
                style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),
              Card(
                elevation: 4,
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'UID: ${_currentUser.uid}',
                        style: const TextStyle(fontWeight: FontWeight.bold),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Nama: ${_currentUser.displayName ?? "Belum diatur"}',
                      ),
                      const SizedBox(height: 8),
                      Text('Email: ${_currentUser.email ?? "Tidak ada email"}'),
                      const SizedBox(height: 8),
                      Text(
                        'Email Terverifikasi: ${_currentUser.emailVerified ? "Ya" : "Tidak"}',
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Terakhir Login: ${_currentUser.metadata.lastSignInTime?.toString() ?? "-"}',
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
