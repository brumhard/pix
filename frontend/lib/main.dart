import 'package:flutter/material.dart';
import 'package:flutter_web_plugins/url_strategy.dart';
import 'package:go_router/go_router.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

void main() {
  usePathUrlStrategy();
  runApp(const SlideShowApp());
}

final _router = GoRouter(
  routes: [
    GoRoute(
      path: '/',
      builder: (context, state) {
        return SlideShow(
          transitionSecondsStr: state.queryParams['transition'] ?? "500",
          intervalSecondsStr: state.queryParams['delay'] ?? "10",
        );
      },
    ),
  ],
  routerNeglect: true,
);

class SlideShowApp extends StatelessWidget {
  const SlideShowApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      routerConfig: _router,
    );
  }
}

class SlideShow extends StatelessWidget {
  final String intervalSecondsStr;
  final String transitionSecondsStr;
  const SlideShow({
    super.key,
    required this.intervalSecondsStr,
    required this.transitionSecondsStr,
  });

  @override
  Widget build(BuildContext context) {
    final channel = WebSocketChannel.connect(
      Uri.parse(
          'ws://${Uri.base.host}:${Uri.base.port}/api/socket?delay=$intervalSecondsStr'),
    );

    return Container(
      decoration: const BoxDecoration(color: Color(0xFF121212)),
      child: StreamBuilder(
        stream: channel.stream,
        builder: (context, snapshot) {
          if (!snapshot.hasData) {
            return Center(
                child: Text(
              "Connecting...",
              style: Theme.of(context).textTheme.headlineLarge?.copyWith(
                    color: Colors.white,
                  ),
            ));
          }

          return AnimatedSwitcher(
            duration: Duration(seconds: int.parse(transitionSecondsStr)),
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
