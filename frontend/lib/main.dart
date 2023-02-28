import 'package:flutter/material.dart';
import 'package:flutter_web_plugins/url_strategy.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

void main() {
  usePathUrlStrategy();
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
    final query = Uri.base.queryParameters;
    final interval = query["delay"] ?? "5";
    final animationStr = query["animation"] ?? "500";
    final animation = int.parse(animationStr);
    final channel = WebSocketChannel.connect(
      Uri.parse(
          'ws://${Uri.base.host}:${Uri.base.port}/api/socket?delay=$interval'),
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
            duration: Duration(milliseconds: animation),
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
