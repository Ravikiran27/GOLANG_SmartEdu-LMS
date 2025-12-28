// TypeScript interfaces matching the Go models

// User types
export interface User {
  uid: string;
  email: string;
  displayName: string;
  role: 'admin' | 'teacher' | 'student';
  photoURL?: string;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
  metadata: UserMetadata;
}

export interface UserMetadata {
  lastLogin?: Date;
  department?: string;
  rollNumber?: string;
  employeeId?: string;
}

export interface CreateUserRequest {
  email: string;
  password: string;
  displayName: string;
  role: 'admin' | 'teacher' | 'student';
  department?: string;
  rollNumber?: string;
  employeeId?: string;
}

export interface UpdateUserRequest {
  displayName?: string;
  photoURL?: string;
  department?: string;
  rollNumber?: string;
  employeeId?: string;
}

// Course types
export interface Course {
  courseId: string;
  title: string;
  description: string;
  syllabus: string;
  teacherId: string;
  teacherName: string;
  category: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  thumbnail?: string;
  materials: CourseMaterial[];
  enrollmentCount: number;
  isPublished: boolean;
  createdAt: Date;
  updatedAt: Date;
  isDeleted: boolean;
}

export interface CourseMaterial {
  id: string;
  name: string;
  type: 'pdf' | 'ppt' | 'video' | 'doc';
  url: string;
  size: number;
  uploadedAt: Date;
}

export interface CreateCourseRequest {
  title: string;
  description: string;
  syllabus?: string;
  category: string;
  difficulty: string;
  thumbnail?: string;
}

export interface UpdateCourseRequest {
  title?: string;
  description?: string;
  syllabus?: string;
  category?: string;
  difficulty?: string;
  thumbnail?: string;
  materials?: CourseMaterial[];
  isPublished?: boolean;
}

export interface Enrollment {
  enrollmentId: string;
  studentId: string;
  studentName: string;
  courseId: string;
  courseTitle: string;
  enrolledAt: Date;
  progress: number;
  completedMaterials: string[];
  status: 'active' | 'completed' | 'dropped';
  lastAccessedAt: Date;
}

// Quiz types
export interface Quiz {
  id: string;
  quizId?: string;
  courseId: string;
  courseTitle?: string;
  teacherId: string;
  title: string;
  description: string;
  duration: number;
  totalMarks: number;
  passingMarks: number;
  negativeMarking: boolean;
  negativeMarkValue: number;
  questionCount: number;
  shuffleQuestions: boolean;
  shuffleOptions: boolean;
  showResultsAfterSubmit: boolean;
  allowReview: boolean;
  maxAttempts: number;
  instructions: string;
  deadline?: Date;
  preventTabSwitch: boolean;
  maxTabSwitches: number;
  requireFullscreen: boolean;
  disableCopyPaste: boolean;
  enableProctoring: boolean;
  randomizeQuestionOrder: boolean;
  timePerQuestion: number;
  lockAfterSubmit: boolean;
  allowTeacherResume: boolean;
  allowTeacherExtendTime: boolean;
  isPublished: boolean;
  createdAt: Date;
  updatedAt: Date;
  isDeleted: boolean;
}

export interface Question {
  id: string;
  questionId?: string;
  quizId?: string;
  type: 'mcq' | 'true_false' | 'short_answer' | 'descriptive';
  text: string;
  questionText?: string;
  imageUrl?: string;
  marks: number;
  points?: number;
  options?: QuestionOption[];
  correctAnswer?: string;
  explanation?: string;
  order: number;
  createdAt: Date;
  updatedAt: Date;
}

export interface QuestionOption {
  id: string;
  text: string;
  isCorrect: boolean;
}

export interface CreateQuizRequest {
  courseId: string;
  title: string;
  description?: string;
  duration: number;
  totalMarks: number;
  passingMarks: number;
  negativeMarking?: boolean;
  negativeMarkValue?: number;
  instructions?: string;
  deadline?: string;
  shuffleQuestions?: boolean;
  shuffleOptions?: boolean;
  showResultsAfterSubmit?: boolean;
  allowReview?: boolean;
  maxAttempts?: number;
  preventTabSwitch?: boolean;
  maxTabSwitches?: number;
  requireFullscreen?: boolean;
  disableCopyPaste?: boolean;
  enableProctoring?: boolean;
  randomizeQuestionOrder?: boolean;
  timePerQuestion?: number;
  allowTeacherResume?: boolean;
  allowTeacherExtendTime?: boolean;
  isPublished?: boolean;
}

export interface CreateQuestionRequest {
  quizId: string;
  type: 'mcq' | 'true_false' | 'short_answer' | 'descriptive';
  questionText: string;
  imageUrl?: string;
  marks: number;
  options?: QuestionOption[];
  correctAnswer?: string;
  explanation?: string;
}

export interface QuizSubmission {
  id: string;
  submissionId?: string;
  quizId: string;
  studentId: string;
  studentName: string;
  studentEmail?: string;
  courseId: string;
  attemptNumber: number;
  answers: Answer[];
  questions?: Question[];
  startedAt: Date;
  submittedAt?: Date;
  timeTaken?: number;
  timeLimit: number;
  totalMarks: number;
  marksObtained: number;
  score?: number;
  percentage: number;
  passed: boolean;
  status: 'in_progress' | 'submitted' | 'evaluated';
  tabSwitchCount: number;
  fullscreenExits: number;
  suspiciousActivity: string[];
  resumedBy?: string;
  resumedAt?: Date;
  resumeReason?: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface Answer {
  questionId: string;
  selectedOptions?: string[];
  textAnswer?: string;
  isCorrect?: boolean;
  pointsAwarded?: number;
}

export interface SubmitQuizRequest {
  submissionId: string;
  quizId: string;
  answers: Answer[];
  tabSwitches?: number;
  fullscreenExits?: number;
  timedOut?: boolean;
}
