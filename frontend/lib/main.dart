import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

void main() {
  runApp(const SlideShowApp());
}

class SlideShowApp extends StatelessWidget {
  const SlideShowApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return const MaterialApp(
      home: SlideShow(),
    );
  }
}

class SlideShow extends StatelessWidget {
  const SlideShow({super.key});

  @override
  Widget build(BuildContext context) {
    final channel = WebSocketChannel.connect(
      Uri.parse('ws://localhost:8888/api/socket'),
    );

    return Container(
      decoration: const BoxDecoration(color: Color(0xFF121212)),
      child: StreamBuilder(
        stream: channel.stream,
        builder: (context, snapshot) {
          if (!snapshot.hasData) {
            return const Text("waiting");
          }

          return AnimatedSwitcher(
            duration: const Duration(milliseconds: 500),
            child: Image.memory(
              snapshot.data,
              // key is needed for the animatedswitcher to get that sth changed
              key: ValueKey(snapshot.data),
              fit: BoxFit.contain,
              height: double.infinity,
              width: double.infinity,
            ),
          );
        },
      ),
    );
  }
}
