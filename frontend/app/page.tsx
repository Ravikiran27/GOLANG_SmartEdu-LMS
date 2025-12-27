export default function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-2xl font-bold text-blue-600">LMS Platform</h1>
            </div>
            <div className="flex space-x-4">
              <a
                href="/auth/login"
                className="px-4 py-2 text-blue-600 hover:text-blue-800"
              >
                Login
              </a>
              <a
                href="/auth/register"
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                Sign Up
              </a>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center">
          <h2 className="text-5xl font-extrabold text-gray-900 mb-6">
            Welcome to Your Learning Journey
          </h2>
          <p className="text-xl text-gray-600 mb-12 max-w-2xl mx-auto">
            A comprehensive Learning Management System with powerful quiz, exam,
            and assignment features designed for modern education.
          </p>

          <div className="grid md:grid-cols-3 gap-8 mt-16">
            <div className="bg-white p-6 rounded-xl shadow-md">
              <div className="text-4xl mb-4">📚</div>
              <h3 className="text-xl font-semibold mb-2">Course Management</h3>
              <p className="text-gray-600">
                Create and manage courses with rich multimedia content
              </p>
            </div>

            <div className="bg-white p-6 rounded-xl shadow-md">
              <div className="text-4xl mb-4">✏️</div>
              <h3 className="text-xl font-semibold mb-2">Quizzes & Exams</h3>
              <p className="text-gray-600">
                Conduct assessments with auto-grading and detailed analytics
              </p>
            </div>

            <div className="bg-white p-6 rounded-xl shadow-md">
              <div className="text-4xl mb-4">📊</div>
              <h3 className="text-xl font-semibold mb-2">Performance Tracking</h3>
              <p className="text-gray-600">
                Monitor student progress with comprehensive analytics
              </p>
            </div>
          </div>

          <div className="mt-16">
            <a
              href="/auth/register"
              className="inline-block px-8 py-4 bg-blue-600 text-white text-lg font-semibold rounded-lg hover:bg-blue-700 shadow-lg transition"
            >
              Get Started Now
            </a>
          </div>
        </div>
      </main>

      <footer className="bg-white border-t mt-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <p className="text-center text-gray-500">
            © 2025 LMS Platform. Built with Next.js, Golang & Firebase.
          </p>
        </div>
      </footer>
    </div>
  );
}
